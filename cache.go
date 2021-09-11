package webcache

import (
	"net/http"
	"sync"
	"time"
)

type item struct {
	request   *http.Response
	touchedAt time.Time
	createdAt time.Time
}

// Cache represents a cache of responses in order to provided it automatically fot a given pattern of request's parameters
type Cache struct {
	mu       sync.RWMutex
	entries  map[string]*item
	timeout  time.Duration
	capacity int
}

func NewCache(capacity int, timeout time.Duration) *Cache {
	return &Cache{
		timeout:  timeout,
		capacity: capacity,
		entries:  make(map[string]*item),
	}
}

func (cache *Cache) release() {
	var target string
	var latest *item

	for key, it := range cache.entries {
		if latest == nil || latest.createdAt.Before(it.createdAt) {
			target = key
			latest = it
		}
	}

	delete(cache.entries, target)
}

func (cache *Cache) GetResponse(tag string) (*http.Response, bool) {
	cache.mu.RLock()
	defer cache.mu.RUnlock()

	item, exists := cache.entries[tag]
	if !exists {
		return nil, false
	}

	item.touchedAt = time.Now()
	if cache.timeout < time.Since(item.createdAt) {
		return nil, false
	}

	return item.request, true
}

func (cache *Cache) SetResponse(tag string, rp *http.Response) {
	cache.mu.Lock()
	defer cache.mu.Unlock()

	if len(cache.entries) >= int(cache.capacity) {
		cache.release()
	}

	cache.entries[tag] = &item{
		request:   rp,
		touchedAt: time.Now(),
		createdAt: time.Now(),
	}
}
