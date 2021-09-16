package webcache

import (
	"log"
	"net/http"
)

// Middleware represents a gateway for getting requests responses
type Middleware interface {
	DecorateRequest(req *http.Request) (*http.Request, error)
	PerformRequest(req *http.Request) (resp *http.Response, err error)
}

// NewHandler returns a brand new request handler for a given middleware
func NewHandler(middle Middleware) http.HandlerFunc {
	return func(wr http.ResponseWriter, req *http.Request) {
		req, err := middle.DecorateRequest(req)
		if err != nil {
			log.Println(err)
			wr.WriteHeader(400)
			wr.Write([]byte(err.Error()))
			return
		}

		if resp, err := middle.PerformRequest(req); err == nil {
			ForwardResponse(resp, wr)
		} else {
			log.Println(err)
			wr.WriteHeader(500)
			wr.Write([]byte(err.Error()))
		}
	}
}
