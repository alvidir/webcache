package webcache

import (
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"sort"
	"sync"
	"time"
)

const (
	ETAG_SERVER_HEADER = "ETag"
	HTTP_CODE_BOUNDARY = 400
	CONCAT_SEPARATOR   = "::"
)

// A Cache represents an storage for request's responses
type Cache interface {
	Store(string, string, time.Duration) error
	Load(string) (string, error)
}

// A FileManager represents an object that creates and removes files
type FileManager interface {
	CreateFile(string) (*os.File, error)
	ReadFile(string) ([]byte, error)
	RemoveFile(string) error
}

// A Manager represents a set of settings about how the ReverseProxy must perform
type Manager interface {
	FileManager

	IsEndpointAllowed(string) bool
	IsMethodAllowed(string, string) bool
	IsMethodCached(string, string) bool
	ResponseLifetime(string) time.Duration
	Headers(string, string) map[string]string
}

// DigestRequest returns the md5 of the given request rq taking as input parameters the request's method,
// the exact host and path, all those listed query params and headers, and the body, if any
func DigestRequest(rq *http.Request, params []string, headers []string) []byte {
	h := md5.New()
	io.WriteString(h, rq.Method)
	io.WriteString(h, rq.Host)

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

	return h.Sum(nil)
}

// A ReverseProxy is a cached reverse proxy that captures responses in order to provide it in the future instead of
// permorming the request each time
type ReverseProxy struct {
	TargetURI       func(req *http.Request) (string, error)
	DigestRequest   func(req *http.Request) (string, error)
	DecorateRequest func(req *http.Request)
	proxys          sync.Map
	manager         Manager
	cache           Cache
}

// NewReverseProxy returns a brand new ReverseProxy with the provided config and cache
func NewReverseProxy(manager Manager, cache Cache) *ReverseProxy {
	reverse := &ReverseProxy{
		manager: manager,
		cache:   cache,
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
	log.Printf("[%s] REQ_TAG %s - %s", req.Method, host, tag)
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

func (reverse *ReverseProxy) getCachedResponseBody(host string, req *http.Request) ([]byte, error) {
	if !reverse.manager.IsMethodCached(host, req.Method) {
		log.Printf("[%s] CACHE_MISS %s", req.Method, host)
		return nil, ErrNotCached
	}

	tag, err := reverse.buildTag(host, req)
	if err != nil {
		return nil, err
	}

	bodyStr, err := reverse.cache.Load(tag)
	if err != nil && err != ErrNotCached {
		log.Printf("[%s] CACHE_MISS %s - %s", req.Method, host, err.Error())
		return nil, err
	}

	if len(bodyStr) == 0 {
		log.Printf("[%s] CACHE_MISS %s", req.Method, host)
		return nil, ErrNoContent
	}

	body, err := base64.RawStdEncoding.DecodeString(bodyStr)
	if err != nil {
		log.Printf("[%s] CACHE_MISS %s - %s", req.Method, host, err.Error())
		return nil, err
	}

	log.Printf("[%s] CACHE_HIT %s", req.Method, host)
	return body, nil
}

func (reverse *ReverseProxy) includeCustomHeaders(host string, req *http.Request) {
	for key, value := range reverse.manager.Headers(host, req.Method) {
		req.Header.Add(key, value)
	}
}

func (reverse *ReverseProxy) storeFileContent(tag string, host string, file *os.File) {
	content, err := ioutil.ReadAll(file)
	if err != nil {
		log.Printf("READ_FILE %s - %s", file.Name(), err.Error())
		return
	}

	log.Printf("STROING: %s", content)
	body := base64.RawStdEncoding.EncodeToString(content)
	timeout := reverse.manager.ResponseLifetime(host)

	if err := reverse.cache.Store(tag, body, timeout); err != nil {
		log.Printf("CACHE_STORE %s - %s", tag, err.Error())
	}
}

func (reverse *ReverseProxy) storeOnSuccess(middleware *HttpMiddleware, tag string, host string, file *os.File) {
	if middleware.header >= HTTP_CODE_BOUNDARY {
		return
	}

	reverse.storeFileContent(tag, host, file)
}

func (reverse *ReverseProxy) performHttpRequest(host string, w http.ResponseWriter, req *http.Request) error {
	proxy, err := reverse.getSingleHostReverseProxy(host)
	if err != nil {
		log.Printf("[%s] PROXY %s - %s", req.Method, host, err.Error())
		return err
	}

	proxy.ErrorLog = log.Default()
	reverse.includeCustomHeaders(host, req)

	tag, err := reverse.buildTag(host, req)
	if err != nil {
		return nil
	}

	if reverse.manager.IsMethodCached(host, req.Method) {
		filename := base64.RawStdEncoding.EncodeToString([]byte(tag))
		file, err := reverse.manager.CreateFile(filename)
		if err != nil {
			return err
		}

		//defer reverse.manager.RemoveFile(filename)
		defer file.Close()

		middleware := NewHttpMiddleware(w, file)
		defer reverse.storeOnSuccess(middleware, tag, host, file)

		proxy.ServeHTTP(middleware, req)
	} else {
		proxy.ServeHTTP(w, req)
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

	if reverse.DecorateRequest != nil {
		reverse.DecorateRequest(req)
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
