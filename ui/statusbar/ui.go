package statusbar

import (
	"fmt"
	"github.com/blackmann/gurl/handler"
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

type Model struct {
	handler  *handler.RequestHandler
	spinner  spinner.Model
	spinning bool
}

func NewStatusBar(handler *handler.RequestHandler) Model {
	s := spinner.New()
	s.Spinner = spinner.Dot

	return Model{handler: handler, spinner: s}
}

func (bar Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	if bar.handler.Status == handler.PROCESSING {
		if !bar.spinning {
			bar.spinning = true
			return bar, bar.spinner.Tick
		}

		var cmd tea.Cmd
		bar.spinner, cmd = bar.spinner.Update(msg)
		return bar, cmd
	} else {

	}

	return bar, nil
}

func (bar Model) View() string {
	switch bar.handler.Status {
	case handler.PROCESSING:
		return fmt.Sprintf("%s Processing", bar.spinner.View())
	default:
		return barStyle.Render(idleStatusStyle.Render("Idle"))
	}
}
