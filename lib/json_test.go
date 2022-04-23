package lib

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseJson_Number(t *testing.T) {
	intTestCases := map[string]int64{"4": 4, "300": 300, "-50": -50}
	for k, n := range intTestCases {
		node, _ := ParseJson([]byte(k), 0)
		assert.Equal(t, n, node.Value[0])
	}

	floatTestCases := map[string]float64{"4.05": 4.05, "0.3": 0.3, "-3.50": -3.50}
	for k, n := range floatTestCases {
		node, _ := ParseJson([]byte(k), 0)
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
		node, _ := ParseJson([]byte(k), 0)
		assert.Equal(t, n, node.Value[0])
	}
}

func TestParseJson_Boolean(t *testing.T) {
	intTestCases := map[string]bool{
		"true":  true,
		"false": false,
	}
	for k, n := range intTestCases {
		node, _ := ParseJson([]byte(k), 0)
		assert.Equal(t, n, node.Value[0])
	}
}

func TestParseJson_Null(t *testing.T) {
	node, _ := ParseJson([]byte("null"), 0)

	assert.Equal(t, node.Type, Null)
}

func TestParseJson_Array(t *testing.T) {
	intTestCases := map[string][]interface{}{
		`["hello", 1]`: {"hello", 10},
	}

	for k, n := range intTestCases {
		node, _ := ParseJson([]byte(k), 0)

		t.Log("> array:", node)

		for i, v := range node.Value {
			t.Log("> array:", n[i], v.(JsonNode).Value)
			assert.Equal(t, v.(JsonNode).Value, n[i])
		}
	}
}
