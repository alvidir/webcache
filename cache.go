package webcache

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/cache/v8"
	"github.com/go-redis/redis/v8"
)

var (
	ErrNotCached = fmt.Errorf("is not cached")
	ErrNoContent = fmt.Errorf("has no content")
)

type RedisCache struct {
	cache *cache.Cache
}

// NewCache returns a basic Cache implementation
func NewCache(size int, ttl time.Duration) (*RedisCache, error) {
	cache := &RedisCache{
		cache: cache.New(&cache.Options{
			LocalCache: cache.NewTinyLFU(size, ttl),
		}),
	}

	return cache, nil
}

// NewRedisCache returns an implementation of Cache for RedisCache
func NewRedisCache(addr string, size int, ttl time.Duration) (*RedisCache, error) {
	ring := redis.NewRing(&redis.RingOptions{
		Addrs: map[string]string{
			"server": addr,
		},
	})

	cache := &RedisCache{
		cache: cache.New(&cache.Options{
			Redis:      ring,
			LocalCache: cache.NewTinyLFU(size, ttl),
		}),
	}

	return cache, nil
}

// Store stores a value under a given key with a lifetime of ttl duration so far
func (c *RedisCache) Store(key string, value any, ttl time.Duration) error {
	ctx := context.Background()
	return c.cache.Set(&cache.Item{
		Ctx:   ctx,
		Key:   key,
		Value: value,
		TTL:   ttl,
	})
}

// Load returns the value for a given key, if any, otherwise err != nil
func (c *RedisCache) Load(key string, value any) (err error) {
	ctx := context.Background()
	if err = c.cache.Get(ctx, key, value); errors.Is(err, cache.ErrCacheMiss) {
		return ErrNotCached
	}

	return
}
