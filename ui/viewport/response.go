package viewport

import (
	"github.com/blackmann/gurl/common"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

type responseModel struct {
	content  string
	viewport viewport.Model
}

func (model responseModel) Update(msg tea.Msg) (responseModel, tea.Cmd) {
	switch msg := msg.(type) {
	case common.Response:
		model.viewport.SetContent(msg.Render())

	case tea.WindowSizeMsg:
		model.viewport.Height = msg.Height
	}

	var cmd tea.Cmd
	model.viewport, cmd = model.viewport.Update(msg)

	return model, cmd
}

func (model responseModel) View() string {
	return model.viewport.View()
}
