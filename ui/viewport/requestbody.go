package viewport

import (
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

type requestBodyModel struct {
	viewport viewport.Model
}

func (model requestBodyModel) Update(msg tea.Msg) (requestBodyModel, tea.Cmd) {
	return model, nil
}

func (model requestBodyModel) View() string {
	return "Request body"
}
