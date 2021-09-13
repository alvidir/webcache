package webcache

import (
	"encoding/base64"
	"log"
	"net/http"
)

// Middleware represents a gateway for getting requests responses
type Middleware interface {
	PerformRequest(req *http.Request) (resp *http.Response, err error)
}

// NewHandler returns a brand new request handler for a given middleware
func NewHandler(mw Middleware) http.HandlerFunc {
	return func(wr http.ResponseWriter, rq *http.Request) {
		query := rq.URL.Query().Get("q")
		uri, err := base64.StdEncoding.DecodeString(query)
		if err != nil {
			log.Println(err)
			wr.WriteHeader(400)
			return
		}

		log.Printf("%s: %s request", uri, rq.Method)
		rq.RequestURI = string(uri)

		if resp, err := mw.PerformRequest(rq); err == nil {
			ForwardResponse(resp, wr)
		} else {
			log.Println(err)
			wr.WriteHeader(500)
			return
		}
	}
}
