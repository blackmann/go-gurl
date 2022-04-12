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
	t.Prompt = "¬ "

	t.Focus()

	return Model{
		input:   t,
		handler: handler,
	}
}

func (model Model) Init() tea.Cmd {
	return textinput.Blink
}

func (model Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			model.handler.MakeRequest()
			return model, nil
		}
	}

	var cmd tea.Cmd
	model.input, cmd = model.input.Update(msg)

	return model, cmd
}

func (model Model) View() string {
	return model.input.View()
}
