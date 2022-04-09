package addressbar

import (
	"github.com/blackmann/gurl/handler"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	input   textinput.Model
	handler *handler.RequestHandler
}

func NewAddressBar(handler *handler.RequestHandler) Model {
	t := textinput.New()
	t.Placeholder = "/GET @adeton/shops"

	t.Focus()

	return Model{
		input:   t,
		handler: handler,
	}
}

func (bar Model) Init() tea.Cmd {
	return textinput.Blink
}

func (bar Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			bar.handler.MakeRequest()
			return bar, nil
		}
	}

	var cmd tea.Cmd
	bar.input, cmd = bar.input.Update(msg)

	return bar, cmd
}

func (bar Model) View() string {
	return bar.input.View()
}
