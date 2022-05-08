package addressbar

import (
	"github.com/blackmann/go-gurl/lib"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestModel_Update_onKeysEntry(t *testing.T) {
	instance := NewAddressBar()

	instance = lib.EnterString(instance, "POST https://example.com").(Model)
	view := instance.View()

	assert.Contains(t, view, "¬ POST https://example.com")
}

func TestModel_Update_onEnter(t *testing.T) {
	instance := NewAddressBar()

	instance = lib.EnterString(instance, "POST https://example.com").(Model)
	_, cmd := instance.Update(tea.KeyMsg{Type: tea.KeyEnter})

	msg := cmd()

	assert.Equal(t, msg, lib.NewRequest)
}

func TestModel_Update_Address(t *testing.T) {
	instance := NewAddressBar()

	tmp, _ := instance.Update(lib.Address{Url: "endpoint", Method: "POST"})
	instance = tmp.(Model)

	addr, _ := instance.(model).GetAddress()
	assert.Equal(t, "endpoint", addr.Url)
	assert.Equal(t, "POST", addr.Method)
}

func TestModel_GetAddress_UrlOnly(t *testing.T) {
	instance := NewAddressBar()
	instance = lib.EnterString(instance, "endpoint").(Model)

	addr, _ := instance.(model).GetAddress()

	assert.Equal(t, "endpoint", addr.Url)
	assert.Equal(t, "GET", addr.Method)
}

func TestModel_GetAddress_UrlAndMethodOnly(t *testing.T) {
	instance := NewAddressBar()
	instance = lib.EnterString(instance, "PATCH endpoint").(Model)

	addr, _ := instance.(model).GetAddress()

	assert.Equal(t, "endpoint", addr.Url)
	assert.Equal(t, "PATCH", addr.Method)
}

func TestModel_GetEntry(t *testing.T) {
	instance := NewAddressBar()
	instance = lib.EnterString(instance, "endpoint").(Model)

	entry := instance.(model).GetEntry()

	assert.Equal(t, "endpoint", entry)
}
