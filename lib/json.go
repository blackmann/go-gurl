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
				// an escape character, so we go ahead and consume the next
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

func ParseJson(str []byte) JsonNode {
	var value []interface{}

	if number, _, err := parseNumber(str, 0); err == nil {
		value = append(value, number)
	}

	if str, _, err := parseString(str, 0); err == nil {
		value = append(value, str)
	}

	return JsonNode{Value: value}
}
