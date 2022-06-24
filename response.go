package webcache

import (
	"fmt"
	"net/http"
	"strings"
)

type HttpResponse struct {
	Body    []byte
	Headers http.Header
	Code    int
}

// implement http.ResponseWriter
// https://golang.org/pkg/net/http/#ResponseWriter
func (r *HttpResponse) Header() http.Header {
	return r.Headers
}

func (r *HttpResponse) Write(b []byte) (int, error) {
	r.Body = append(r.Body, b...)
	return len(r.Body), nil
}

func (r *HttpResponse) WriteHeader(i int) {
	r.Code = i
}

func (r *HttpResponse) Echo(w http.ResponseWriter) (int, error) {
	for header, values := range r.Headers {
		for _, value := range values {
			w.Header().Add(header, value)
		}
	}

	w.WriteHeader(r.Code)
	return w.Write(r.Body)
}

func (r *HttpResponse) Format() (format string) {
	format = fmt.Sprintf("HTTP %v\n", r.Code)
	for header, values := range r.Headers {
		format += fmt.Sprintf("%s: %s\n", header, strings.Join(values, ", "))
	}

	if r.Body == nil || len(r.Body) == 0 {
		format += "EMPTY_BODY"
		return
	}

	format += fmt.Sprintf("BODY_LEN %v bytes", len(r.Body))
	return
}

func (r *HttpResponse) Empty() bool {
	return len(r.Body) == 0 &&
		len(r.Headers) == 0 &&
		r.Code == 0
}

func NewHttpResponse() *HttpResponse {
	return &HttpResponse{
		Body:    []byte{},
		Headers: make(http.Header),
		Code:    0,
	}
}
