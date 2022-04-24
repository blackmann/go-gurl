package lib

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
)

func HighlightJson(content string) (string, error) {
	decoder := json.NewDecoder(strings.NewReader(content))
	decoder.UseNumber()

	return highlight(decoder)
}

func highlight(decoder *json.Decoder) (string, error) {
	for {
		t, err := decoder.Token()

		if err == io.EOF {
			break
		}

		if delim, ok := t.(json.Delim); ok {
			switch delim {
			case '[':
				return highlightArray(decoder), nil
			case '{':
				return highlightObject(decoder), nil
			}
		}

		switch t.(type) {
		case json.Number:
			return highlightNumber(t), nil

		case string:
			return highlightString(t), nil

		case nil:
			return highlightNull(), nil

		case bool:
			return highlightBool(t), nil
		}
	}

	return "", errors.New("invalid json")
}

func highlightNumber(token json.Token) string {
	return fmt.Sprintf("%v", token)
}

func highlightString(token json.Token) string {
	return fmt.Sprintf(`"%v"`, token)
}

func highlightBool(token json.Token) string {
	return fmt.Sprintf("%v", token)
}

func highlightNull() string {
	return "null"
}

func highlightArray(decoder *json.Decoder) string {
	s := "["

	for decoder.More() {
		if entry, err := highlight(decoder); err == nil {
			if len(s) > 1 {
				s += fmt.Sprintf(", \n%s", entry)
			} else {
				s += fmt.Sprintf("\n%s", entry)
			}
		} else {
			// TODO: Handle
		}
	}

	_, _ = decoder.Token() // Clean the trailing ]

	return s + "]"
}

func highlightObject(decoder *json.Decoder) string {
	s := "{"

	for decoder.More() {
		token, _ := decoder.Token()
		key := fmt.Sprintf(`"%v"`, token)
		value, _ := highlight(decoder)

		if len(s) > 1 {
			s += fmt.Sprintf(",\n%s: %s", key, value)
		} else {
			s += fmt.Sprintf("\n%s: %s", key, value)
		}
	}

	return s + "\n}"
}
