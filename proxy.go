package webcache

import (
	"crypto/md5"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sort"
	"sync"
	"time"
)

const (
	ETAG_SERVER_HEADER = "ETag"
	HTTP_CODE_BOUNDARY = 400
	CONCAT_SEPARATOR   = "::"
)

var (
	ErrNotCached = fmt.Errorf("is not cached")
	ErrNoContent = fmt.Errorf("has no content")
)

// A Cache represents an storage for request's responses
type Cache interface {
	Store(string, string, time.Duration) error
	Load(string) (string, error)
}

// A Config represents a set of settings about how the ReverseProxy must perform
type Config interface {
	IsEndpointAllowed(string) bool
	IsMethodAllowed(string, string) bool
	IsMethodCached(string, string) bool
	ResponseLifetime(string) time.Duration
	Headers(string, string) map[string]string
}

type httpMiddleware struct {
	body   []byte
	header int
	w      http.ResponseWriter
}

// implement http.ResponseWriter
// https://golang.org/pkg/net/http/#ResponseWriter
func (cacher *httpMiddleware) Header() http.Header {
	return cacher.w.Header()
}

func (cacher *httpMiddleware) Write(b []byte) (int, error) {
	cacher.body = append(cacher.body, b...)
	return cacher.w.Write(b)
}

func (cacher *httpMiddleware) WriteHeader(i int) {
	cacher.header = i
	cacher.w.WriteHeader(i)
}

func newHttpMiddleware(w http.ResponseWriter) *httpMiddleware {
	return &httpMiddleware{
		body: make([]byte, 1024),
		w:    w,
	}
}

// DigestRequest returns the md5 of the given request rq taking as input parameters the request's method,
// the exact host and path, all those listed query params and headers, and the body, if any
func DigestRequest(rq *http.Request, params []string, headers []string) string {
	h := md5.New()
	io.WriteString(h, rq.Method)
	io.WriteString(h, rq.RequestURI)

	sort.Strings(params)
	for _, param := range params {
		label := fmt.Sprintf("%s%s", param, rq.URL.Query().Get(param))
		io.WriteString(h, label)
	}

	sort.Strings(headers)
	for _, header := range headers {
		label := fmt.Sprintf("%s%s", header, rq.Header.Get(header))
		io.WriteString(h, label)
	}

	if body, err := io.ReadAll(rq.Body); err == nil {
		io.WriteString(h, string(body))
	}

	return string(h.Sum(nil))
}

// A ReverseProxy is a cached reverse proxy that captures responses in order to provide it in the future instead of
// permorming the request each time
type ReverseProxy struct {
	TargetURI     func(req *http.Request) (string, error)
	DigestRequest func(req *http.Request) (string, error)
	proxys        sync.Map
	config        Config
	cache         Cache
}

// NewReverseProxy returns a brand new ReverseProxy with the provided config and cache
func NewReverseProxy(config Config, cache Cache) *ReverseProxy {
	reverse := &ReverseProxy{
		config: config,
		cache:  cache,
	}

	return reverse
}

func (reverse *ReverseProxy) buildTag(host string, req *http.Request) (tag string, err error) {
	tag = req.Header.Get(ETAG_SERVER_HEADER)
	if reverse.DigestRequest != nil {
		tag, err = reverse.DigestRequest(req)
		if err != nil {
			log.Printf("[%s] REQ_TAG %s - %s", req.Method, host, err.Error())
			return "", err
		}
	}

	tag = fmt.Sprintf("%s::%s::%s", req.Method, host, tag)
	return
}

func (reverse *ReverseProxy) getSingleHostReverseProxy(host string) (*httputil.ReverseProxy, error) {
	if v, ok := reverse.proxys.Load(host); ok {
		if fn, ok := v.(*httputil.ReverseProxy); ok {
			return fn, nil
		}

		log.Printf("TYPE_ASSERT %s - want *httputil.ReverseProxy", host)
	}

	remoteUrl, err := url.Parse(host)
	if err != nil {
		log.Printf("URL_PARSE %s - %s", host, err.Error())
		return nil, err
	}

	proxy := httputil.NewSingleHostReverseProxy(remoteUrl)
	reverse.proxys.Store(host, proxy)
	return proxy, nil
}

func (reverse *ReverseProxy) getCachedResponseBody(host string, req *http.Request) (string, error) {
	if !reverse.config.IsMethodCached(host, req.Method) {
		log.Printf("[%s] CACHE_MISS %s", req.Method, host)
		return "", ErrNotCached
	}

	tag, err := reverse.buildTag(host, req)
	if err != nil {
		return "", err
	}

	body, err := reverse.cache.Load(tag)
	if err != nil {
		log.Printf("[%s] CACHE_MISS %s - %s", req.Method, host, err.Error())
		return "", err
	}

	if len(body) == 0 {
		log.Printf("[%s] CACHE_MISS %s", req.Method, host)
		return "", ErrNoContent
	}

	log.Printf("[%s] CACHE_HIT %s", req.Method, host)
	return body, nil
}

func (reverse *ReverseProxy) includeCustomHeaders(host string, req *http.Request) {
	for key, value := range reverse.config.Headers(host, req.Method) {
		req.Header.Add(key, value)
	}
}

func (reverse *ReverseProxy) performHttpRequest(host string, w http.ResponseWriter, req *http.Request) error {
	proxy, err := reverse.getSingleHostReverseProxy(host)
	if err != nil {
		log.Printf("[%s] PROXY %s - %s", req.Method, host, err.Error())
		return err
	}

	reverse.includeCustomHeaders(host, req)

	middleware := newHttpMiddleware(w)
	proxy.ServeHTTP(middleware, req)

	if !reverse.config.IsMethodCached(host, req.Method) {
		return nil
	}

	tag, err := reverse.buildTag(host, req)
	if err != nil {
		return nil
	}

	if middleware.header < HTTP_CODE_BOUNDARY {
		timeout := reverse.config.ResponseLifetime(host)
		reverse.cache.Store(tag, string(middleware.body), timeout)
	}

	return nil
}

// ServeHTTP performs http requests if not cached yet or returns the chaced body instead
func (reverse *ReverseProxy) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	host := req.Host
	if reverse.TargetURI != nil {
		var err error
		if host, err = reverse.TargetURI(req); err != nil {
			log.Printf("TARGET_URI %s", err)

			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("400: Bad request"))
			return
		}
	}

	if !reverse.config.IsEndpointAllowed(host) {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("403: Host forbidden " + host))
		return
	}

	if !reverse.config.IsMethodAllowed(host, req.Method) {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("405: Method not allowed " + host))
		return
	}

	if body, err := reverse.getCachedResponseBody(host, req); err == nil {
		w.Write([]byte(body))
		return
	} else if err != ErrNotCached && err != ErrNoContent {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500: Internal server error"))
		return
	}

	if err := reverse.performHttpRequest(host, w, req); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500: Internal server error"))
	}
}
