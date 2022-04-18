package statusbar

import (
	"fmt"
	"github.com/blackmann/gurl/lib"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// A message sent with a value to set as the status' value
type statusUpdate lib.Status

type Model struct {
	spinner  spinner.Model
	spinning bool

	width        int
	status       lib.Status
	commandEntry string
}

func NewStatusBar() Model {
	s := spinner.New()
	s.Spinner = spinner.Dot

	return Model{spinner: s}
}

func (model Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case commandInput:
		model.commandEntry = string(msg)
		return model, nil

	case statusUpdate:
		model.status = lib.Status(msg)
		//return model, nil
	}

	if model.status == lib.PROCESSING {
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
	case lib.PROCESSING:
		view = fmt.Sprintf("%s Processing", model.spinner.View())
	default:
		view = barStyle.Render(idleStatusStyle.Render("Idle"))
	}

	if len(model.commandEntry) > 0 {
		view = fmt.Sprintf("%s %s", view, model.commandEntry)
	}

	return lipgloss.NewStyle().Margin(1).Render(view)
}
