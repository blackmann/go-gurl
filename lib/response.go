package lib

import (
	"net/http"
	"strings"
)

type Response struct {
	Body    []byte
	Headers http.Header
	Status  int
	Time    int64
}

func (response Response) Render() string {
	content := ""

	contentType := response.Headers.Get("content-type")

	if strings.HasPrefix(contentType, "application/json") {
		p := ColoredPrettier()
		if pretty, err := p.HighlightJson(string(response.Body)); err == nil {
			content = pretty
		} else {
			content = string(response.Body)
		}
	}

	return content
}
