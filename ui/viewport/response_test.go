package viewport

import (
	"errors"
	"github.com/blackmann/go-gurl/lib"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestResponseModel_Init(t *testing.T) {
	// In a running application, initialization actually happens
	// when tea.WindowMsg is passed
	model, _ := responseModel{}.Update(nil)

	assert.Contains(t, model.View(), "No response data")
}

func TestResponseModel_Update_Response(t *testing.T) {
	model, _ := responseModel{}.Update(tea.WindowSizeMsg{Width: 30, Height: 50})

	model, _ = model.Update(lib.Response{Headers: http.Header{}, Body: []byte("raw text")})

	assert.Contains(t, model.View(), "raw text")
}

func TestResponseModel_Update_RequestError(t *testing.T) {
	model, _ := responseModel{}.Update(tea.WindowSizeMsg{Width: 30, Height: 50})

	model, _ = model.Update(lib.RequestError{Err: errors.New("Invalid url schema")})

	assert.Contains(t, model.View(), "Invalid url schema")
}
