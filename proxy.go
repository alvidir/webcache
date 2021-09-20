package webcache

import "sync"

// import "net/url"

// var reverseProxy = NewSingleHostReverseProxy(target*url.URL) * ReverseProxy

// Conn represents a db connection
type Conn interface {
	Store(key, value string)
	Load(key string) (string, bool)
}

type ReverseProxy struct {
	configByEnpoint  sync.Map
	configByFilename sync.Map
	db               Conn
}

func NewReverseProxy(dbConn Conn) *ReverseProxy {
	return &ReverseProxy{
		db: dbConn,
	}
}
