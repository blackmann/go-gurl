package lib

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHighlightJson_Number_Int(t *testing.T) {
	v, _ := HighlightJson(`20`)
	assert.Equal(t, "20", v)
}

func TestHighlightJson_Number_Float(t *testing.T) {
	v, _ := HighlightJson(`20.50`)
	assert.Equal(t, "20.50", v)
}

func TestHighlightJson_String(t *testing.T) {
	v, _ := HighlightJson(`"hello"`)
	assert.Equal(t, `"hello"`, v)
}

func TestHighlightJson_Bool_True(t *testing.T) {
	v, _ := HighlightJson(`true`)
	assert.Equal(t, "true", v)
}

func TestHighlightJson_Bool_False(t *testing.T) {
	v, _ := HighlightJson(`false`)
	assert.Equal(t, "false", v)
}

func TestHighlightJson_Null(t *testing.T) {
	v, _ := HighlightJson(`null`)
	assert.Equal(t, "null", v)
}

func TestHighlightJson_Array(t *testing.T) {
	testCases := map[string]string{
		`[1]`:                     "[1]",
		`[1, 3, false]`:           "[1, 3, false]",
		`[1, null]`:               "[1, null]",
		`[1, [1, 40.5], "hello"]`: "[1, [1, 40.5], \"hello\"]",
		`[{}]`:                    "[{}]",
	}

	for input, expectation := range testCases {
		v, _ := HighlightJson(input)
		assert.Equal(t, expectation, v)
	}
}

func TestHighlightJson_Object(t *testing.T) {
	v, _ := HighlightJson(`{"key": 1}`)
	assert.Equal(t, "{}", v)
}
