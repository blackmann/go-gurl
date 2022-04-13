package addressbar

import (
	"github.com/blackmann/gurl/common/commands"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	input textinput.Model
}

func NewAddressBar() Model {
	t := textinput.New()
	t.Placeholder = "/GET @adeton/shops"
	t.Prompt = "¬ "

	t.Focus()

	return Model{
		input: t,
	}
}

func (model Model) Init() tea.Cmd {
	textinput.Blink()
	return textinput.Blink
}

func (model Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			return model, commands.SubmitNewRequest
		}
	}

	var cmd tea.Cmd
	model.input, cmd = model.input.Update(msg)

	return model, cmd
}

func (model Model) View() string {
	return model.input.View()
}
