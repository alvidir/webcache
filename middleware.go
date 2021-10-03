package webcache

import (
	"io"
	"net/http"
)

type HttpMiddleware struct {
	file    io.Writer
	resp    http.ResponseWriter
	multi   io.Writer
	headers http.Header
	code    int
}

// implement http.ResponseWriter
// https://golang.org/pkg/net/http/#ResponseWriter
func (w *HttpMiddleware) Header() http.Header {
	return w.resp.Header()
}

func (w *HttpMiddleware) Write(b []byte) (int, error) {
	return w.multi.Write(b)
}

func (w *HttpMiddleware) WriteHeader(i int) {
	w.resp.WriteHeader(i)

	w.code = i
	w.headers = w.resp.Header()
}

func NewHttpMiddleware(resp http.ResponseWriter, file io.Writer) *HttpMiddleware {
	multi := io.MultiWriter(resp, file)
	return &HttpMiddleware{
		file:    file,
		resp:    resp,
		multi:   multi,
		headers: make(http.Header),
		code:    0,
	}
}
