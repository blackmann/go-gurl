package viewport

import (
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

type headersModel struct {
	viewport viewport.Model
}

func (model headersModel) Update(msg tea.Msg) (headersModel, tea.Cmd) {
	return model, nil
}

func (model headersModel) View() string {
	return "Headers"
}
