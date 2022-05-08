package viewport

import (
	"github.com/blackmann/go-gurl/lib"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestHeadersModel_Init(t *testing.T) {
	model, _ := headersModel{}.Update(nil)

	// The cursor appears at the beginning of the placeholder
	// let's rely on the rest of the string (without H)
	assert.Contains(t, model.View(), "eader-Key: Value")
}

func TestHeadersModel_Update_Header(t *testing.T) {
	model, _ := headersModel{}.Update(nil)

	header := http.Header{}
	header.Add("accept", "application/json")
	model, _ = model.Update(header)

	view := model.View()
	// Gets canonicalized
	assert.Contains(t, view, "Accept")
	assert.Contains(t, view, "application/json")
}

func TestHeadersModel_Update_HeaderInput(t *testing.T) {
	model, _ := headersModel{}.Update(nil)
	header := "accept: application/json"
	model = lib.EnterString(model, header)

	assert.Contains(t, model.View(), header)

	var cmd tea.Cmd
	model, cmd = model.Update(tea.KeyMsg{Type: tea.KeyEnter})

	pair := cmd().(lib.Pair)

	assert.Equal(t, pair.Key, "accept")
	assert.Equal(t, pair.Value, "application/json")

	assert.NotContains(t, model.View(), header)
}

func TestHeadersModel_Update_NavigateList(t *testing.T) {
	model, _ := headersModel{}.Update(nil)

	header := http.Header{}
	header.Add("accept", "application/json")
	header.Add("authorization", "bearer 123")
	model, _ = model.Update(header)

	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyDown})
	assert.Equal(t, model.(headersModel).verticalPosition, LIST)

	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyDown})
	assert.Equal(t, model.(headersModel).list.SelectedItem().(lib.Pair).Key, "Accept")

	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyUp})
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyUp})

	assert.Equal(t, model.(headersModel).verticalPosition, INPUT)
}

func TestHeadersModel_Update_InputIntoList(t *testing.T) {
	// when in the list position, don't receive inputs for the
	// header input

	model, _ := headersModel{}.Update(nil)

	header := http.Header{}
	header.Add("accept", "application/json")
	model, _ = model.Update(header)

	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyDown})
	assert.Equal(t, model.(headersModel).verticalPosition, LIST)

	model = lib.EnterString(model, "new-entry: ")

	assert.NotContains(t, model.View(), "new-entry")
}
