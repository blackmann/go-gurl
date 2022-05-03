package viewport

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
	"testing"
)

// FIX: This is copy-pasted everywhere, is there a generic solution?
func enterBodyString(m requestBodyModel, keys string) requestBodyModel {
	for _, k := range keys {
		hKey := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{k}}
		m, _ = m.Update(hKey)
	}

	return m
}

func TestRequestBody_Update_InputEntry(t *testing.T) {
	model := requestBodyModel{}

	model = enterBodyString(model, "name: Mock")
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	model = enterBodyString(model, "age: 23")

	assert.Contains(t, model.View(), "  name: Mock  \n  age: 23")
}

func TestRequestBody_Update_RequestBody(t *testing.T) {
	model := requestBodyModel{}

	model, _ = model.Update(requestBody("hello\nworld"))
	assert.Contains(t, model.View(), "  hello   \n  world")
}
