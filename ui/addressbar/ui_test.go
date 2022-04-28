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

func TestModel_Update_Address(t *testing.T) {
	instance := NewAddressBar()

	instance, _ = instance.Update(lib.Address{Url: "endpoint", Method: "POST"})

	addr, _ := instance.GetAddress()
	assert.Equal(t, "endpoint", addr.Url)
	assert.Equal(t, "POST", addr.Method)
}

func TestModel_GetAddress_UrlOnly(t *testing.T) {
	instance := NewAddressBar()
	instance = enterString(instance, "endpoint")

	addr, _ := instance.GetAddress()

	assert.Equal(t, "endpoint", addr.Url)
	assert.Equal(t, "GET", addr.Method)
}

func TestModel_GetAddress_UrlAndMethodOnly(t *testing.T) {
	instance := NewAddressBar()
	instance = enterString(instance, "PATCH endpoint")

	addr, _ := instance.GetAddress()

	assert.Equal(t, "endpoint", addr.Url)
	assert.Equal(t, "PATCH", addr.Method)
}

func TestModel_GetEntry(t *testing.T) {
	instance := NewAddressBar()
	instance = enterString(instance, "endpoint")

	entry := instance.GetEntry()

	assert.Equal(t, "endpoint", entry)
}
