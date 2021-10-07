package webcache

import (
	"fmt"
	"net/http"
	"strings"
)

type HttpResponse struct {
	body    []byte
	headers http.Header
	code    int
}

// implement http.ResponseWriter
// https://golang.org/pkg/net/http/#ResponseWriter
func (r *HttpResponse) Header() http.Header {
	return r.headers
}

func (r *HttpResponse) Write(b []byte) (int, error) {
	r.body = append(r.body, b...)
	return len(r.body), nil
}

func (r *HttpResponse) WriteHeader(i int) {
	r.code = i
}

func (r *HttpResponse) Echo(w http.ResponseWriter) {
	for header, values := range r.headers {
		for _, value := range values {
			w.Header().Add(header, value)
		}
	}

	w.WriteHeader(r.code)
	w.Write(r.body)
}

func (r *HttpResponse) Cat() {
	fmt.Printf("HTTP %v", r.code)
	for header, values := range r.headers {
		fmt.Printf("%s: %s", header, strings.Join(values, ", "))
	}

	fmt.Printf("\n %s", string(r.body))
}

func NewHttpResponse() HttpResponse {
	return HttpResponse{
		body:    []byte{},
		headers: make(http.Header),
		code:    0,
	}
}
