package webcache

import (
	"crypto/md5"
	"io"
	"net/http"
	"sort"
)

func HashRequest(rq *http.Request, headers []string) string {
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

func ForwardResponse(resp *http.Response, wr http.ResponseWriter) (err error) {
	defer resp.Body.Close()
	wr.WriteHeader(resp.StatusCode)

	for key, values := range resp.Header {
		for _, value := range values {
			wr.Header().Add(key, value)
		}
	}

	_, err = io.Copy(wr, resp.Body)
	return
}
