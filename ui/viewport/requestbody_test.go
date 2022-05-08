package viewport

import (
	"github.com/blackmann/go-gurl/lib"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRequestBody_Update_InputEntry(t *testing.T) {
	model := newRequestBodyModel()

	model = lib.EnterString(model, "name: Mock").(RequestBodyModel)
	tmp, _ := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	model = tmp.(RequestBodyModel)

	model = lib.EnterString(model, "age: 23").(RequestBodyModel)

	assert.Contains(t, model.View(), "  name: Mock  \n  age: 23")
}

func TestRequestBody_Update_RequestBody(t *testing.T) {
	model := newRequestBodyModel()

	tmp, _ := model.Update(requestBody("hello\nworld"))
	model = tmp.(RequestBodyModel)
	assert.Contains(t, model.View(), "  hello   \n  world")
}
