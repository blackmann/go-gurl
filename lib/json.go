package lib

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
)

type NodeType int

var (
	Number  NodeType = 1
	String  NodeType = 2
	Boolean NodeType = 3
	Null    NodeType = 4
	Array   NodeType = 5
	Object  NodeType = 6
)

type JsonNode struct {
	Type  NodeType
	Value []interface{}
	Key   string
}

var numberRegex = regexp.MustCompile(`\d`)

func parseNumber(str []byte, cursor int) (interface{}, int, error) {
	s := ""
	decimalFound := false

	for ; cursor < len(str); cursor++ {
		char := str[cursor]
		if numberRegex.Match([]byte{char}) {
			s += string(char)
		} else if char == '.' && len(s) > 0 && !decimalFound {
			s += "."
			decimalFound = true
		} else if char == '-' && len(s) == 0 {
			s += "-"
		} else {
			break
		}
	}

	if len(s) == 0 {
		return nil, cursor, errors.New("not a number")
	}

	if decimalFound {
		if number, err := strconv.ParseFloat(s, 64); err != nil {
			return nil, cursor, err
		} else {
			return number, cursor, nil
		}
	}

	if number, err := strconv.ParseInt(s, 10, 64); err != nil {
		return nil, cursor, err
	} else {
		return number, cursor, nil
	}
}

func parseString(str []byte, cursor int) (string, int, error) {
	s := ""

	if str[cursor] != '"' {
		return "", cursor, errors.New("not a string")
	}
	cursor++

	for ; cursor < len(str); cursor++ {
		char := str[cursor]
		if char != '"' {
			if char != '\\' {
				s += string(char)
			} else {
				// an escape character, so we go ahead and consume the next \"..."
				cursor++
				s += fmt.Sprintf("%c%c", char, str[cursor])
			}
		} else {
			cursor++
			break
		}
	}

	return s, cursor, nil
}

func parseBool(str []byte, cursor int) (bool, int, error) {
	if str[cursor] == 't' {
		// true, advance cursor
		cursor += 4
		return true, cursor, nil
	}

	if str[cursor] == 'f' {
		cursor += 5
		return false, cursor, nil
	}

	return false, cursor, errors.New("not a bool")
}

func parseNull(str []byte, cursor int) (interface{}, int, error) {
	if str[cursor] == 'n' {
		return nil, cursor + 4, nil
	}
	return nil, cursor, errors.New("not a null")
}

func parseArray(str []byte, cursor int) ([]interface{}, int, error) {
	if str[cursor] == '[' {
		cursor++
		var res []interface{}
		for str[cursor] != ']' {
			if str[cursor] == ',' || str[cursor] == ' ' {
				cursor++
				continue
			}

			if node, _cursor, err := ParseJson(str, cursor); err == nil {
				res = append(res, node)
				cursor = _cursor
			} else {
				return nil, cursor, err
			}
		}

		cursor++ // close the "]"
		return res, cursor, nil
	}

	return nil, cursor, errors.New("not an array")
}

func ParseJson(str []byte, cursor int) (JsonNode, int, error) {
	if number, _cursor, err := parseNumber(str, cursor); err == nil {
		return JsonNode{Value: []interface{}{number}, Type: Number}, _cursor, nil
	}

	if str, _cursor, err := parseString(str, cursor); err == nil {
		return JsonNode{Value: []interface{}{str}, Type: String}, _cursor, nil
	}

	if boolean, _cursor, err := parseBool(str, cursor); err == nil {
		return JsonNode{Value: []interface{}{boolean}, Type: Boolean}, _cursor, nil
	}

	if _, _cursor, err := parseNull(str, cursor); err == nil {
		return JsonNode{Type: Null, Value: []interface{}{nil}}, _cursor, nil
	}

	if arr, _cursor, err := parseArray(str, cursor); err == nil {
		return JsonNode{Type: Array, Value: arr}, _cursor, nil
	}

	return JsonNode{}, cursor, errors.New("invalid json")
}
