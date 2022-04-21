package lib

import (
	"bytes"
	"encoding/json"
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
		var pretty bytes.Buffer
		if err := json.Indent(&pretty, response.Body, "", "  "); err == nil {
			content = pretty.String()
		} else {
			content = string(response.Body)
		}
	}

	return content
}
