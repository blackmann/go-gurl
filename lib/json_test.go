package lib

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseJson_Number(t *testing.T) {
	intTestCases := map[string]int64{"4": 4, "300": 300, "-50": -50}
	for k, n := range intTestCases {
		node := ParseJson([]byte(k))
		assert.Equal(t, n, node.Value[0])
	}

	floatTestCases := map[string]float64{"4.05": 4.05, "0.3": 0.3, "-3.50": -3.50}
	for k, n := range floatTestCases {
		node := ParseJson([]byte(k))
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
		node := ParseJson([]byte(k))
		assert.Equal(t, n, node.Value[0])
	}
}
