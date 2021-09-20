package webcache

import (
	"github.com/go-redis/redis"
)

type rclient struct {
	*redis.Client
}

// NewRedisClient returns an implementation of Conn for a redis db
func NewRedisConn(url string) (Conn, error) {
	opt, err := redis.ParseURL(url)
	if err != nil {
		return nil, err
	}

	client := rclient{
		redis.NewClient(opt),
	}

	return client, nil
}

func (client rclient) Store(key, value string) {

}

func (client rclient) Load(key string) (string, bool) {
	return "", false
}
