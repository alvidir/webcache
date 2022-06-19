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
	"strings"
	"sync"
	"time"
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

// A Manager represents a set of settings about how the ReverseProxy must perform
type Manager interface {
	IsEndpointAllowed(string) bool
	IsMethodAllowed(string, string) bool
	IsMethodCached(string, string) bool
	ResponseLifetime(string) time.Duration
	Headers(string, string) map[string]string
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
	manager       Manager
	responses     Cache
}

// NewReverseProxy returns a brand new ReverseProxy with the provided config and cache
func NewReverseProxy(manager Manager, cache Cache) *ReverseProxy {
	reverse := &ReverseProxy{
		manager:   manager,
		responses: cache,
	}

	return reverse
}

func (reverse *ReverseProxy) targetURI(req *http.Request) (host string, err error) {
	if targets, ok := req.Header[HTTP_LOCATION_HEADER]; ok && len(targets) > 0 {
		return targets[0], nil
	}

	return "", ErrNoContent
}

func (reverse *ReverseProxy) tag(host string, req *http.Request) (tag string, err error) {
	tag = req.Header.Get(ETAG_SERVER_HEADER)
	if reverse.DigestRequest != nil {
		tag, err = reverse.DigestRequest(req)
		if err != nil {
			log.Printf("[%s] REQ_TAG %s - %s", req.Method, host, err.Error())
			return
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

	proxy := &httputil.ReverseProxy{
		Director: func(req *http.Request) {
			req.Header.Add("X-Forwarded-Host", req.Host)
			req.Header.Add("X-Target-Host", remoteUrl.Host)
			req.URL.Scheme = remoteUrl.Scheme
			req.URL.Host = remoteUrl.Host
			req.Host = remoteUrl.Host

			req.Header.Del(HTTP_LOCATION_HEADER)
			//log.Printf("REQUEST\n%s\n", FormatHttpRequest(req))
		},
	}

	reverse.proxys.Store(host, proxy)
	return proxy, nil
}

func (reverse *ReverseProxy) getCachedResponseBody(host string, req *http.Request) (*HttpResponse, error) {
	if !reverse.manager.IsMethodCached(host, req.Method) {
		log.Printf("[%s] CACHE_MISS %s", req.Method, host)
		return nil, ErrNotCached
	}

	tag, err := reverse.tag(host, req)
	if err != nil {
		return nil, err
	}

	resp := NewHttpResponse()
	if err := reverse.responses.Load(tag, resp); err != nil && err != ErrNotCached {
		log.Printf("[%s] CACHE_MISS %s - %s", req.Method, tag, err.Error())
		return nil, ErrNotCached
	} else if resp.Empty() {
		log.Printf("[%s] CACHE_MISS %s - %s", req.Method, tag, ErrNotCached.Error())
		return nil, ErrNotCached
	}

	log.Printf("[%s] CACHE_HIT %s - %+v", req.Method, tag, resp)
	return resp, nil
}

func (reverse *ReverseProxy) includeCustomHeaders(host string, req *http.Request) {
	for key, value := range reverse.manager.Headers(host, req.Method) {
		req.Header.Add(key, value)
	}
}

func (reverse *ReverseProxy) performHttpRequest(w http.ResponseWriter, req *http.Request, host string) error {
	proxy, err := reverse.getSingleHostReverseProxy(host)
	if err != nil {
		log.Printf("[%s] PROXY %s - %s", req.Method, host, err.Error())
		return err
	}

	proxy.ErrorLog = log.Default()
	reverse.includeCustomHeaders(host, req)

	if !reverse.manager.IsMethodCached(host, req.Method) {
		proxy.ServeHTTP(w, req)
		return nil
	}

	response := NewHttpResponse()
	proxy.ServeHTTP(response, req)

	//log.Printf("RESPONSE\n%s\n", response.Format())

	if diff := response.Code - HTTP_CODE_REDIRECT; 0 <= diff && diff < 100 {
		// as HTTP_CODE_REDIRECT == 300, then diff is somewhere between 300 and 399
		headers := response.Header()
		if locations, exists := headers[HTTP_LOCATION_HEADER]; !exists || len(locations) == 0 {
			log.Printf("REDIRECT %s - %s http header must be set", host, HTTP_LOCATION_HEADER)
			return nil
		} else if len(locations) > 1 {
			log.Printf("REDIRECT %s - %s http header has too much values", host, HTTP_LOCATION_HEADER)
			return nil
		}

		location := headers[HTTP_LOCATION_HEADER][0]
		location = strings.Split(location, "?")[0]
		if location == host {
			log.Printf("REDIRECT %s - cyclical redirection to itself", host)
			response.Echo(w)
			return nil
		}

		log.Printf("REDIRECT %s - has been moved to %s", host, location)
		return reverse.performHttpRequest(w, req, location)
	}

	if response.Code < HTTP_CODE_BOUNDARY {
		go func() {
			tag, err := reverse.tag(host, req)
			if err != nil {
				return
			}

			timeout := reverse.manager.ResponseLifetime(host)
			if err := reverse.responses.Store(tag, response, timeout); err != nil {
				log.Printf("CACHE_STORE %s - %s", tag, err.Error())
			} else {
				log.Printf("CACHE_STORE %s", tag)
			}
		}()
	}

	response.Echo(w)
	return nil
}

// ServeHTTP performs http requests if not cached yet or returns the chaced body instead
func (reverse *ReverseProxy) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	host, err := reverse.targetURI(req)
	if err != nil {
		log.Printf("TARGET_URI %s", err)

		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("400: Bad request"))
		return
	}

	if !reverse.manager.IsEndpointAllowed(host) {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("403: Host forbidden " + host))
		return
	}

	if !reverse.manager.IsMethodAllowed(host, req.Method) {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("405: Method not allowed " + host))
		return
	}

	if resp, err := reverse.getCachedResponseBody(host, req); err == nil {
		resp.Echo(w)
		return
	} else if err != ErrNotCached && err != ErrNoContent {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500: Internal server error"))
		return
	}

	if err := reverse.performHttpRequest(w, req, host); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500: Internal server error"))
	}
}
