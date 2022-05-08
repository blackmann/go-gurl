package viewport

import (
	"github.com/blackmann/go-gurl/lib"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type responseModel struct {
	initialized bool
	viewport    viewport.Model

	hasResponse bool
	height      int
	width       int
}

func (model *responseModel) initialize() {
	// Doing this to maintain a sturdy viewport.
	// Removing this causes the viewport to jiggle because
	// there's no content
	model.viewport.SetContent("")

	// This only works in tea.WithMouseCellMotion()
	//model.viewport.MouseWheelEnabled = true

	model.initialized = true
}

func (model responseModel) Init() tea.Cmd {
	return nil
}

func (model responseModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if !model.initialized {
		model.initialize()
	}

	contentStyle := lipgloss.NewStyle().Width(model.viewport.Width)

	switch msg := msg.(type) {
	case lib.RequestError:
		contentStyle = contentStyle.Foreground(lipgloss.Color(lib.ANSIRed))
		model.viewport.SetContent(contentStyle.Render(msg.Err.Error()))
		model.viewport.YOffset = 0

		model.hasResponse = true

	case lib.Response:
		model.viewport.SetContent(contentStyle.Render(msg.Render(true)))
		model.viewport.YOffset = 0

		model.hasResponse = true

	case tea.WindowSizeMsg:
		model.viewport.Height = msg.Height
		model.viewport.Width = msg.Width - 4 // Left/right padding

		model.height = msg.Height
		model.width = msg.Width

	case lib.Event:
		if msg == lib.Reset {
			model.viewport.SetContent("")
			model.hasResponse = false
		}
	}

	var cmd tea.Cmd
	model.viewport, cmd = model.viewport.Update(msg)

	return model, cmd
}

func (model responseModel) View() string {
	if !model.hasResponse {
		style := lipgloss.NewStyle().Foreground(lipgloss.Color("#999")).
			Padding(0, 2, 0, 2).Height(model.height)
		return style.Render("No response data")
	}

	return lipgloss.NewStyle().Padding(0, 2, 0, 2).Render(model.viewport.View())
}
