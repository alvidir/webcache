package webcache

import (
	"crypto/md5"
	"io"
	"net/http"
	"sync"
	"time"
)

// Cache responses for a given key
const (
	MISS uint = iota
	HIT
)

type item struct {
	v         http.Request
	touchedAt time.Time
	createdAt time.Time
}

// Cache represents a cache of responses in order to provided it automatically fot a given pattern of request's parameters
type Cache struct {
	mu      sync.RWMutex
	entries map[string]item
}

func getRequestChecksum(rq *http.Request) string {
	h := md5.New()
	io.WriteString(h, rq.Method)
	io.WriteString(h, rq.RequestURI)

	return string(h.Sum(nil))
}
