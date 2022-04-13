package statusbar

import (
	"fmt"
	"github.com/blackmann/gurl/common/status"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	barStyle = lipgloss.NewStyle().
			Padding(0, 1, 0, 1)

	idleStatusStyle = lipgloss.NewStyle().Background(lipgloss.Color("#fff")).
			Foreground(lipgloss.Color("#000")).
			Padding(0, 1, 0, 1)
)

type updateCommand string

type widthUpdate int

type Model struct {
	spinner  spinner.Model
	spinning bool

	width        int
	status       status.Status
	commandEntry string
}

func Command(command string) tea.Msg {
	return updateCommand(command)
}

func Width(width int) tea.Msg {
	return widthUpdate(width)
}

func NewStatusBar() Model {
	s := spinner.New()
	s.Spinner = spinner.Dot

	return Model{spinner: s}
}

func (model Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case updateCommand:
		model.commandEntry = string(msg)
		return model, nil
	}

	if model.status == status.PROCESSING {
		if !model.spinning {
			model.spinning = true
			return model, model.spinner.Tick
		}

		var cmd tea.Cmd
		model.spinner, cmd = model.spinner.Update(msg)
		return model, cmd
	} else {
		model.spinning = false
	}

	return model, nil
}

func (model Model) View() string {
	var view string

	switch model.status {
	case status.PROCESSING:
		view = fmt.Sprintf("%s Processing", model.spinner.View())
	default:
		view = barStyle.Render(idleStatusStyle.Render("Idle"))
	}

	if len(model.commandEntry) > 0 {
		view = fmt.Sprintf("%s %s", view, model.commandEntry)
	}

	return lipgloss.NewStyle().Margin(1).Render(view)
}

func (model Model) SetStatus(status status.Status) Model {
	model.status = status
	return model
}
