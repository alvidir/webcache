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

func (r *HttpResponse) Format() (format string) {
	format = fmt.Sprintf("HTTP %v\n", r.code)
	for header, values := range r.headers {
		format += fmt.Sprintf("%s: %s\n", header, strings.Join(values, ", "))
	}

	format += fmt.Sprintf("\n%s\n", string(r.body))
	return
}

func NewHttpResponse() *HttpResponse {
	return &HttpResponse{
		body:    []byte{},
		headers: make(http.Header),
		code:    0,
	}
}
