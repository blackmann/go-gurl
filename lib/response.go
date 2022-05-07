package lib

import (
	"bytes"
	"encoding/json"
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

func (response Response) Render(pretty bool) string {
	contentType := response.Headers.Get("content-type")

	if strings.HasPrefix(contentType, "application/json") {
		if pretty {
			p := ColoredPrettier()
			if prettyJson, err := p.HighlightJson(string(response.Body)); err == nil {
				return prettyJson
			}
		}

		var out bytes.Buffer
		json.Indent(&out, response.Body, "", "  ")

		return out.String()
	}

	return string(response.Body)
}
