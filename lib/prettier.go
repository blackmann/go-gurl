package lib

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"io"
	"strings"
)

type palette struct {
	null   lipgloss.Style
	bool   lipgloss.Style
	number lipgloss.Style
	string lipgloss.Style
	delim  lipgloss.Style
	key    lipgloss.Style
}

type Prettier struct {
	palette
}

func ColoredPrettier() Prettier {
	colorPalette := palette{
		null:   lipgloss.NewStyle(),
		bool:   lipgloss.NewStyle().Foreground(lipgloss.Color(ANSICyan)),
		number: lipgloss.NewStyle().Foreground(lipgloss.Color(ANSIBlue)),
		string: lipgloss.NewStyle().Foreground(lipgloss.Color(ANSIGreen)),
		delim:  lipgloss.NewStyle(),
		key:    lipgloss.NewStyle().Foreground(lipgloss.Color(ANSIRed)),
	}

	return Prettier{palette: colorPalette}
}

func NoColorPrettier() Prettier {
	return Prettier{}
}

func getIndent(depth int) string {
	return fmt.Sprintf("\n%s", strings.Repeat("  ", depth))
}

func (p Prettier) HighlightJson(content string) (string, error) {
	decoder := json.NewDecoder(strings.NewReader(content))
	decoder.UseNumber()

	return p.highlight(decoder, 0)
}

func (p Prettier) highlight(decoder *json.Decoder, depth int) (string, error) {
	for {
		t, err := decoder.Token()

		if err == io.EOF {
			break
		}

		if delim, ok := t.(json.Delim); ok {
			switch delim {
			case '[':
				return p.highlightArray(decoder, depth), nil
			case '{':
				return p.highlightObject(decoder, depth), nil
			}
		}

		switch t.(type) {
		case json.Number:
			return p.highlightNumber(t), nil

		case string:
			return p.highlightString(t), nil

		case nil:
			return p.highlightNull(), nil

		case bool:
			return p.highlightBool(t), nil
		}
	}

	return "", errors.New("invalid json")
}

func (p Prettier) highlightNumber(token json.Token) string {
	return p.palette.number.Render(fmt.Sprintf("%v", token))
}

func (p Prettier) highlightString(token json.Token) string {
	return p.palette.string.Render(fmt.Sprintf(`"%v"`, token))
}

func (p Prettier) highlightBool(token json.Token) string {
	return p.palette.bool.Render(fmt.Sprintf("%v", token))
}

func (p Prettier) highlightNull() string {
	return p.palette.null.Render("null")
}

func (p Prettier) highlightArray(decoder *json.Decoder, depth int) string {
	s := p.palette.delim.Render("[")

	isEmpty := true
	innerIndent := getIndent(depth + 1)

	for decoder.More() {
		isEmpty = false
		if entry, err := p.highlight(decoder, depth+1); err == nil {
			if len(s) > 1 {
				s += fmt.Sprintf(",%s%s", innerIndent, entry)
			} else {
				s += fmt.Sprintf("%s%s", innerIndent, entry)
			}
		} else {
			// TODO: Handle
		}
	}

	_, _ = decoder.Token() // Clean the trailing: ]

	closingDelim := p.palette.delim.Render("]")
	if isEmpty {
		return s + closingDelim
	}

	closingIndent := getIndent(depth)
	return fmt.Sprintf("%s%s%s", s, closingIndent, closingDelim)
}

func (p Prettier) highlightObject(decoder *json.Decoder, depth int) string {
	s := p.palette.delim.Render("{")
	isEmpty := true
	innerIndent := getIndent(depth + 1)

	for decoder.More() {
		isEmpty = false
		token, _ := decoder.Token()
		key := p.palette.key.Render(fmt.Sprintf(`"%v"`, token))
		value, _ := p.highlight(decoder, depth+1)

		if len(s) > 1 {
			s += fmt.Sprintf(",%s%s: %s", innerIndent, key, value)
		} else {
			s += fmt.Sprintf("%s%s: %s", innerIndent, key, value)
		}
	}

	_, _ = decoder.Token() // Clean the trailing: }

	closingDelim := p.palette.delim.Render("}")
	if isEmpty {
		return s + closingDelim
	}

	closingIndent := getIndent(depth)
	return fmt.Sprintf("%s%s%s", s, closingIndent, closingDelim)
}
