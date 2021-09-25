package webcache

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
	"time"
)

var (
	ErrNotCached = fmt.Errorf("is not cached")
	ErrNoContent = fmt.Errorf("has no content")
)

// Cache represents an storage for request's responses
type Cache interface {
	Store(string, string, time.Duration)
	Load(string) (string, error)
}

type Config interface {
	IsMethodAllowed(string, string) bool
	IsMethodCached(string, string) bool
}

type ReverseProxy struct {
	TargetURI   func(req *http.Request) (string, error)
	HashRequest func(req *http.Request) (string, error)
	proxys      sync.Map
	config      Config
	cache       Cache
}

func NewReverseProxy(config Config, cache Cache) *ReverseProxy {
	reverse := &ReverseProxy{
		config: config,
		cache:  cache,
	}

	return reverse
}

func (reverse *ReverseProxy) getSingleHostReverseProxy(host string) (*httputil.ReverseProxy, error) {
	if v, ok := reverse.proxys.Load(host); ok {
		if fn, ok := v.(*httputil.ReverseProxy); ok {
			return fn, nil
		}

		log.Printf("got wrong type, want *httputil.ReverseProxy")
	}

	remoteUrl, err := url.Parse(host)
	if err != nil {
		log.Println("target parse fail:", err)
		return nil, err
	}

	proxy := httputil.NewSingleHostReverseProxy(remoteUrl)
	reverse.proxys.Store(host, proxy)
	return proxy, nil
}

func (reverse *ReverseProxy) getCachedResponseBody(host string, req *http.Request) (string, error) {
	if !reverse.config.IsMethodCached(host, req.Method) ||
		reverse.HashRequest == nil {
		return "", ErrNotCached
	}

	key, err := reverse.HashRequest(req)
	if err != nil {
		return "", err
	}

	body, err := reverse.cache.Load(key)
	if err != nil {
		return "", err
	}

	if len(body) == 0 {
		return "", ErrNoContent
	}

	return body, nil
}

func (reverse *ReverseProxy) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	host := req.Host
	if reverse.TargetURI != nil {
		var err error
		if host, err = reverse.TargetURI(req); err != nil {
			w.Write([]byte("400: Bad request"))
		}
	}

	if !reverse.config.IsMethodAllowed(host, req.Method) {
		w.Write([]byte("403: Host forbidden " + host))
	}

	if body, err := reverse.getCachedResponseBody(host, req); err == nil {
		log.Printf("[%s] CACHE_HIT %s", req.Method, host)
		w.Write([]byte(body))
		return
	} else if err != ErrNotCached && err != ErrNoContent {
		log.Printf("[%s] CACHE_MISS %s - %s", req.Method, host, err.Error())
		w.Write([]byte("500: Internal server error"))
		return
	}

	log.Printf("[%s] CACHE_MISS %s", req.Method, host)
	proxy, err := reverse.getSingleHostReverseProxy(host)
	if err != nil {
		log.Printf("[%s] PROXY_ERROR %s - %s", req.Method, host, err.Error())
		w.Write([]byte("500: Internal server error"))
		return
	}

	proxy.ServeHTTP(w, req)
}
