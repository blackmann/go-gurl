package lib

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestHighlightJson_Number_Int(t *testing.T) {
	p := ColoredPrettier()
	v, _ := p.HighlightJson(`20`)
	assert.Equal(t, "20", v)
}

func TestHighlightJson_Number_Float(t *testing.T) {
	p := ColoredPrettier()
	v, _ := p.HighlightJson(`20.50`)
	assert.Equal(t, "20.50", v)
}

func TestHighlightJson_String(t *testing.T) {
	p := ColoredPrettier()
	v, _ := p.HighlightJson(`"hello"`)
	assert.Equal(t, `"hello"`, v)
}

func TestHighlightJson_Bool_True(t *testing.T) {
	p := ColoredPrettier()
	v, _ := p.HighlightJson(`true`)
	assert.Equal(t, "true", v)
}

func TestHighlightJson_Bool_False(t *testing.T) {
	p := ColoredPrettier()
	v, _ := p.HighlightJson(`false`)
	assert.Equal(t, "false", v)
}

func TestHighlightJson_Null(t *testing.T) {
	p := ColoredPrettier()
	v, _ := p.HighlightJson(`null`)
	assert.Equal(t, "null", v)
}

func TestHighlightJson_Array(t *testing.T) {
	p := ColoredPrettier()
	testCases := map[string]string{
		// Single
		`[1]`: strings.Join([]string{
			`[`,
			`  1`,
			`]`,
		}, "\n"),
		// Multiple, mixed
		`[1, "hello", false, null]`: strings.Join([]string{
			`[`,
			`  1,`,
			`  "hello",`,
			`  false,`,
			`  null`,
			`]`,
		}, "\n"),
		// Nested
		`[1, [true, 40.5], "hello", {"key": "value"}]`: strings.Join([]string{
			`[`,
			`  1,`,
			`  [`,
			`    true,`,
			`    40.5`,
			`  ],`,
			`  "hello",`,
			`  {`,
			`    "key": "value"`,
			`  }`,
			`]`,
		}, "\n"),
		// Empty
		`[]`: "[]",
	}

	for input, expectation := range testCases {
		v, _ := p.HighlightJson(input)
		assert.Equal(t, expectation, v)
	}
}

func TestHighlightJson_Object(t *testing.T) {
	input := `{"8": "eight", "__v": {"package": "dev", "version": ["2019.03", 22]}}`

	expectation := strings.Join([]string{
		`{`,
		`  "8": "eight",`,
		`  "__v": {`,
		`    "package": "dev",`,
		`    "version": [`,
		`      "2019.03",`,
		`      22`,
		`    ]`,
		`  }`,
		`}`,
	}, "\n")

	p := ColoredPrettier()

	v, _ := p.HighlightJson(input)
	assert.Equal(t, expectation, v)
}
