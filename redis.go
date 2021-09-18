package webcache

import (
	"crypto/tls"

	"github.com/go-redis/redis"
)

var redisClient = redis.NewClient(&redis.Options{
	Network:            "",
	Addr:               "redis:6379",
	Dialer:             nil, // func() (net.Conn, error)
	OnConnect:          nil, // func(*redis.Conn) error
	Password:           "",
	DB:                 0,
	MaxRetries:         0,
	MinRetryBackoff:    0,
	MaxRetryBackoff:    0,
	DialTimeout:        0,
	ReadTimeout:        0,
	WriteTimeout:       0,
	PoolSize:           0,
	MinIdleConns:       0,
	MaxConnAge:         0,
	PoolTimeout:        0,
	IdleTimeout:        0,
	IdleCheckFrequency: 0,
	TLSConfig:          &tls.Config{},
})

func SetRedisOption(options *redis.Options) {
	redisClient = redis.NewClient(options)
}
