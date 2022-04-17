package statusbar

import (
	"fmt"
	"github.com/blackmann/gurl/common/status"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// A command in this case is an entry from the keyboard
// that is mapped to an action. For example, ":q" to quit, etc.
// The action is not necessarily performed by statusbar. We're
// only using this type as a state update message type.
//
//  statusbar.Update(commandInput(":q"))
//
// This is a tea.Msg type
type commandInput string

// A message sent with a value to be set as the statusbar's width
type widthUpdate int

// A message sent with a value to set as the status' value
type statusUpdate status.Status

type Model struct {
	spinner  spinner.Model
	spinning bool

	width        int
	status       status.Status
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
		model.status = status.Status(msg)
		//return model, nil
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
