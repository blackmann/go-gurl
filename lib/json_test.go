package lib

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseJson_Number(t *testing.T) {
	intTestCases := map[string]int64{"4": 4, "300": 300, "-50": -50}
	for k, n := range intTestCases {
		node, _, _ := ParseJson([]byte(k), 0)
		assert.Equal(t, n, node.Value[0])
	}

	floatTestCases := map[string]float64{"4.05": 4.05, "0.3": 0.3, "-3.50": -3.50}
	for k, n := range floatTestCases {
		node, _, _ := ParseJson([]byte(k), 0)
		assert.Equal(t, n, node.Value[0])
	}
}

func TestParseJson_String(t *testing.T) {
	intTestCases := map[string]string{
		`"hello"`:             "hello",
		`"world help"`:        "world help",
		`"world \"apart\""`:   `world \"apart\"`,
		`"world\n \"apart\""`: `world\n \"apart\"`,
	}
	for k, n := range intTestCases {
		node, _, _ := ParseJson([]byte(k), 0)
		assert.Equal(t, n, node.Value[0])
	}
}

func TestParseJson_Boolean(t *testing.T) {
	intTestCases := map[string]bool{
		"true":  true,
		"false": false,
	}
	for k, n := range intTestCases {
		node, _, _ := ParseJson([]byte(k), 0)
		assert.Equal(t, n, node.Value[0])
	}
}

func TestParseJson_Null(t *testing.T) {
	node, _, _ := ParseJson([]byte("null"), 0)

	assert.Equal(t, node.Type, Null)
}

func TestParseJson_Array(t *testing.T) {
	intTestCases := map[string][]interface{}{
		`["hello", 1]`:        {"hello", 1},
		`["hello", null, 50]`: {"hello", nil, 50},
		`["hello", 0, false]`: {"hello", 0, false},
		`["hello", 0, true]`:  {"hello", 0, true},
		`[]`:                  {},
		`[null]`:              {nil},
		// Nesting
		`[["hello", 0, true], 5]`:       {[]interface{}{"hello", 0, true}, 5},
		`[null, ["hello", 0, true], 5]`: {nil, []interface{}{"hello", 0, true}, 5},
	}

	for testCase, testExpectation := range intTestCases {
		node, _, err := ParseJson([]byte(testCase), 0)

		assert.Nil(t, err)

		assert.Equal(t, node.Type, Array)
		assert.Equal(t, len(testExpectation), len(node.Value))

		for i, v := range node.Value {
			expected := testExpectation[i]

			switch expected.(type) {
			// For nested cases
			case []interface{}:
				expectedArray := expected.([]interface{})
				for i2, v2 := range expectedArray {
					assert.EqualValues(t, v2, v.(JsonNode).Value[i2].(JsonNode).Value[0])
				}
			default:
				assert.EqualValues(t, testExpectation[i], v.(JsonNode).Value[0])
			}
		}
	}
}
