package addressbar

import (
	"github.com/blackmann/gurl/common/appcmd"
	"github.com/blackmann/gurl/common/request"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"strings"
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
			return model, appcmd.SubmitNewRequest
		}
	}

	var cmd tea.Cmd
	model.input, cmd = model.input.Update(msg)

	return model, cmd
}

func (model Model) View() string {
	return model.input.View()
}

func (model Model) GetAddress() request.Address {
	return request.Address{Method: "GET", Url: strings.Trim(model.input.Value(), " ")}
}
