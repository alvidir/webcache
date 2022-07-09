package webcache

import (
	"crypto/md5"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sort"
	"strings"
	"sync"
	"time"

	"go.uber.org/zap"
)

const (
	ETAG_SERVER_HEADER   = "ETag"
	HTTP_LOCATION_HEADER = "Location"
	HTTP_CODE_BOUNDARY   = 400
	HTTP_CODE_REDIRECT   = 300
)

// A Cache represents an storage for request's responses
type Cache interface {
	Store(string, interface{}, time.Duration) error
	Load(string, interface{}) error
}

// A Config represents a set of settings about how the ReverseProxy must perform
type Config interface {
	RequestOptions(endpoint, method string) (*Options, bool)
}

// DigestRequest returns the md5 of the given request rq taking as input parameters the request's method,
// the exact host and path, all those listed query params and headers, and the body, if any
func DigestRequest(rq *http.Request, headers []string) []byte {
	h := md5.New()
	io.WriteString(h, rq.Method)
	io.WriteString(h, rq.Host)
	io.WriteString(h, rq.URL.Query().Encode())

	sort.Strings(headers)
	for _, header := range headers {
		label := fmt.Sprintf("%s%s", header, rq.Header.Get(header))
		io.WriteString(h, label)
	}

	if body, err := io.ReadAll(rq.Body); err == nil {
		io.WriteString(h, string(body))
	}

	return h.Sum(nil)
}

// FormatHttpRequest returns and string descriving the content of an HttpRequest
func FormatHttpRequest(req *http.Request) (format string) {
	format = fmt.Sprintf("%s %s\n\n", req.Method, req.URL)
	for header, values := range req.Header {
		format += fmt.Sprintf("%s: %s\n", header, strings.Join(values, ", "))
	}

	if req.Body == nil {
		format += "EMPTY_BODY"
		return
	}

	if bytes, err := io.ReadAll(req.Body); err != nil {
		format += fmt.Sprintf("BODY_ERROR %s", err.Error())
	} else if len(bytes) > 0 {
		format += fmt.Sprintf("BODY_LEN %v bytes", len(bytes))
	} else {
		format += "EMPTY_BODY"
	}

	return
}

// A ReverseProxy is a cached reverse proxy that captures responses in order to provide it in the future instead of
// permorming the same request each time
type ReverseProxy struct {
	DigestRequest func(req *http.Request) (string, error)
	proxys        sync.Map
	manager       Config
	responses     Cache
	logger        *zap.Logger
}

// NewReverseProxy returns a brand new ReverseProxy with the provided config and cache
func NewReverseProxy(manager Config, cache Cache, logger *zap.Logger) *ReverseProxy {
	reverse := &ReverseProxy{
		manager:   manager,
		responses: cache,
		logger:    logger,
	}

	return reverse
}

func (reverse *ReverseProxy) target(req *http.Request) (host string, err error) {
	if targets, ok := req.Header[HTTP_LOCATION_HEADER]; ok && len(targets) > 0 {
		return targets[0], nil
	}

	return "", ErrNoContent
}

func (reverse *ReverseProxy) tag(req *http.Request) (string, error) {
	if reverse.DigestRequest == nil {
		return req.Header.Get(ETAG_SERVER_HEADER), nil
	}

	host, _ := reverse.target(req)
	tag, err := reverse.DigestRequest(req)
	if err != nil {
		reverse.logger.Error("digesting http request",
			zap.String("host", host),
			zap.String("method", req.Method),
			zap.Error(err))

		return "", ErrUnknown
	}

	return fmt.Sprintf("%s::%s::%s", req.Method, host, tag), nil
}

func (reverse *ReverseProxy) singleHostReverseProxy(host string) (*httputil.ReverseProxy, error) {
	if v, ok := reverse.proxys.Load(host); ok {
		return v.(*httputil.ReverseProxy), nil
	}

	remoteUrl, err := url.Parse(host)
	if err != nil {
		reverse.logger.Error("parsing url",
			zap.String("url", host),
			zap.Error(err))

		return nil, ErrUnknown
	}

	proxy := httputil.NewSingleHostReverseProxy(remoteUrl)
	proxy.Director = func(req *http.Request) {
		req.Header.Del(HTTP_LOCATION_HEADER)
		req.URL.Scheme = remoteUrl.Scheme
		req.URL.Host = remoteUrl.Host
		req.Host = remoteUrl.Host
	}

	reverse.proxys.Store(host, proxy)
	return proxy, nil
}

func (reverse *ReverseProxy) storeResponseBody(req *http.Request, resp *HttpResponse, ops *Options) {
	tag, _ := reverse.tag(req)
	timeout := ops.timeout

	if err := reverse.responses.Store(tag, resp, timeout); err != nil {
		reverse.logger.Error("storing response",
			zap.String("tag", tag),
			zap.Error(err))
	} else {
		reverse.logger.Info("response stored",
			zap.String("tag", tag))
	}
}

func (reverse *ReverseProxy) loadResponseBody(req *http.Request, ops *Options) (*HttpResponse, bool) {
	host, _ := reverse.target(req)

	if !ops.cached {
		reverse.logger.Info("cache miss",
			zap.String("host", host),
			zap.String("method", req.Method))

		return nil, false
	}

	tag, err := reverse.tag(req)
	if err != nil {
		return nil, false
	}

	resp := NewHttpResponse()
	if err := reverse.responses.Load(tag, &resp); err != nil &&
		!errors.Is(err, ErrNotFound) {
		reverse.logger.Warn("loading response from cache",
			zap.String("host", host),
			zap.String("method", req.Method),
			zap.Error(err))
	}

	if resp.Empty() {
		reverse.logger.Info("cache miss",
			zap.String("host", host),
			zap.String("method", req.Method))

		return nil, false
	}

	reverse.logger.Warn("cache hit",
		zap.String("host", host),
		zap.String("method", req.Method),
		zap.Error(err))

	return resp, true
}

func (reverse *ReverseProxy) addHeaders(req *http.Request, ops *Options) {
	for key, value := range ops.headers {
		req.Header.Add(key, value)
	}
}

func (reverse *ReverseProxy) follow(req *http.Request, ops *Options, host string) (*HttpResponse, error) {
	reverse.logger.Info("performing request",
		zap.String("host", host))

	proxy, err := reverse.singleHostReverseProxy(host)
	if err != nil {
		return nil, err
	}

	resp := NewHttpResponse()
	proxy.ServeHTTP(resp, req)

	if resp.Code/100 != HTTP_CODE_REDIRECT/100 {
		return resp, nil
	}

	// response.Code is redirect: 3XX
	locs, exists := resp.Header()[HTTP_LOCATION_HEADER]
	if !exists || len(locs) == 0 {
		reverse.logger.Error("getting location header",
			zap.String("host", host),
			zap.Error(ErrNotFound))

		return nil, ErrNotFound
	}

	target := locs[0]
	target = strings.Split(target, "?")[0]
	if target == host {
		reverse.logger.Warn("cyclical redirection",
			zap.String("host", host),
			zap.String("location", target))

		return resp, nil
	}

	reverse.logger.Info("host has been moved",
		zap.String("host", host),
		zap.String("location", target))

	return reverse.follow(req, ops, target)
}

func (reverse *ReverseProxy) perform(w http.ResponseWriter, req *http.Request, ops *Options) {
	host, _ := reverse.target(req)
	resp, err := reverse.follow(req, ops, host)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500: Internal server error"))
	}

	if ops.cached && resp.Code < HTTP_CODE_BOUNDARY {
		reverse.storeResponseBody(req, resp, ops)
	}

	if _, err := resp.Echo(w); err != nil {
		reverse.logger.Error("sending response back to the client",
			zap.Error(err))
	}
}

// ServeHTTP performs http requests if not cached yet or returns the chaced body instead
func (reverse *ReverseProxy) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	host, err := reverse.target(req)
	if err != nil {
		reverse.logger.Error("getting request's target host",
			zap.Error(err))

		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("400: Bad request"))
		return
	}

	ops, _ := reverse.manager.RequestOptions(host, req.Method)
	if ops == nil {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("403: Host forbidden " + host))
		return
	}

	if !ops.enabled {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("405: Method not allowed " + host))
		return
	}

	if resp, exists := reverse.loadResponseBody(req, ops); exists {
		resp.Echo(w)
		return
	}

	reverse.addHeaders(req, ops)
	reverse.perform(w, req, ops)
}
