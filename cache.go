package webcache

import (
	"sync"
	"time"
)

type cacheItem struct {
	response  []byte
	touchedAt time.Time
}

type Cache struct {
	entries  sync.Map
	mu       sync.RWMutex
	capacity uint
	size     uint
}
