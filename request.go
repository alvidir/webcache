package webcache

import (
	"crypto/md5"
	"io"
	"net/http"
	"sort"
)

func RequestChecksum(rq *http.Request, headers []string) string {
	h := md5.New()
	io.WriteString(h, rq.Method)
	io.WriteString(h, rq.RequestURI)

	sort.Strings(headers)
	for _, header := range headers {
		label := header + rq.Header.Get(header)
		io.WriteString(h, label)
	}

	if body, err := io.ReadAll(rq.Body); err == nil {
		io.WriteString(h, string(body))
	}

	return string(h.Sum(nil))
}

func RequestDecorator(req *http.Request, headers map[string]string) {
	for key, value := range headers {
		req.Header.Add(key, value)
	}
}

func PerformRequest(req *http.Request, config *Config) *http.Response {
	return nil
}

func ForwardResponse(resp *http.Response, wr http.ResponseWriter) {
	wr.WriteHeader(resp.StatusCode)

	for key, values := range resp.Header {
		for _, value := range values {
			wr.Header().Add(key, value)
		}
	}

	io.Copy(wr, resp.Body)
	resp.Body.Close()
}
