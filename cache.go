package webcache

import (
	"sync"
	"time"
)

// Cache responses for a given key
const (
	MISS uint = iota
	HIT
)

type cacheItem struct {
	v         interface{}
	touchedAt time.Time
}

// Cache represents a cache of responses in order to provided it automatically fot a given pattern of request's parameters
type Cache struct {
	mu       sync.RWMutex
	entries  map[string]cacheItem
	capacity uint
	size     uint
}
