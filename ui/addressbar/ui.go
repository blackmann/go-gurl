package addressbar

import (
	"errors"
	"fmt"
	"github.com/blackmann/go-gurl/lib"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"strings"
)

type Model interface {
	tea.Model
	GetAddress() (lib.Address, error)
	GetEntry() string
}

type model struct {
	input textinput.Model
}

func (model model) Init() tea.Cmd {
	return textinput.Blink
}

func NewAddressBar() Model {
	t := textinput.New()
	t.Placeholder = "GET @endpoint/path"
	t.Prompt = "¬ "

	t.Focus()

	return model{input: t}
}

func (model model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			submitRequest := func() tea.Msg { return lib.NewRequest }
			return model, submitRequest
		}

	case lib.Address:
		model.input.SetValue(fmt.Sprintf("%s %s", msg.Method, msg.Url))
		model.input.CursorEnd()
		return model, nil

	case string:
		model.input.SetValue(msg)
		model.input.CursorEnd()

		return model, nil
	}

	var cmd tea.Cmd
	model.input, cmd = model.input.Update(msg)

	return model, cmd
}

func (model model) View() string {
	return model.input.View()
}

func (model model) GetAddress() (lib.Address, error) {
	trimmedAddr := strings.Trim(model.input.Value(), " ")

	if len(trimmedAddr) == 0 {
		return lib.Address{}, errors.New("no address entry")
	}

	// Expecting at least one part and most 2
	// When two, the first part is the method and the latter is the endpoint
	parts := strings.Split(trimmedAddr, " ")

	method := "GET"
	endpoint := parts[len(parts)-1]

	if len(parts) > 1 {
		method = parts[0]
	}

	return lib.Address{Method: strings.ToUpper(method), Url: endpoint}, nil
}

func (model model) GetEntry() string {
	return model.input.Value()
}
