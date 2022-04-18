package viewport

import (
	"github.com/blackmann/gurl/lib"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

type responseModel struct {
	initialized bool
	viewport    viewport.Model
}

func (model responseModel) Update(msg tea.Msg) (responseModel, tea.Cmd) {
	if !model.initialized {
		// Doing this to maintain a sturdy viewport.
		// Removing this causes the viewport to jiggle because
		// there's no content
		model.viewport.SetContent("")
		model.initialized = true
	}

	switch msg := msg.(type) {
	case lib.Response:
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
