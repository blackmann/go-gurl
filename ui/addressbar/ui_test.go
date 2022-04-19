package addressbar

import (
	"github.com/blackmann/gurl/lib"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
	"testing"
)

func enterString(m Model, keys string) Model {
	for _, k := range keys {
		hKey := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{k}}
		m, _ = m.Update(hKey)
	}

	return m
}

func TestModel_Update_onKeysEntry(t *testing.T) {
	instance := NewAddressBar()

	instance = enterString(instance, "POST https://example.com")
	view := instance.View()

	assert.Contains(t, view, "¬ POST https://example.com")
}

func TestModel_Update_onEnter(t *testing.T) {
	instance := NewAddressBar()

	instance = enterString(instance, "POST https://example.com")
	_, cmd := instance.Update(tea.KeyMsg{Type: tea.KeyEnter})

	msg := cmd()

	assert.Equal(t, msg, lib.NewRequest)
}
