package lib

import (
	"net/http"
	"strings"
)

type Response struct {
	Body    []byte
	Cookies []*http.Cookie
	Headers http.Header
	Status  int
	Time    int64
	Request Request
}

func (response Response) Render() string {
	contentType := response.Headers.Get("content-type")

	if strings.HasPrefix(contentType, "application/json") {
		p := ColoredPrettier()
		if pretty, err := p.HighlightJson(string(response.Body)); err == nil {
			return pretty
		}
	}

	return string(response.Body)
}
