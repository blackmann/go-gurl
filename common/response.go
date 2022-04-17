package common

import (
	"bytes"
	"encoding/json"
)

type Response struct {
	ContentType string
	Body        []byte
}

func (response Response) Render() string {
	content := ""

	if response.ContentType == "application/json" {
		var pretty bytes.Buffer
		if err := json.Indent(&pretty, response.Body, "", "  "); err == nil {
			content = pretty.String()
		} else {
			content = string(response.Body)
		}
	}

	return content
}
