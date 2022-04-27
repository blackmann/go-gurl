package viewport

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"strings"
)

var encodedEnter = "¬"

type requestBody string

type requestBodyModel struct {
	input textinput.Model

	initialized bool
	height      int
}

func (model requestBodyModel) Input() string {
	return strings.Replace(model.input.Value(), encodedEnter, "\n", -1)
}

func (model *requestBodyModel) Init() tea.Cmd {
	input := textinput.New()
	input.Focus()
	input.Placeholder = ""
	input.Prompt = ""

	model.input = input
	model.initialized = true

	return nil
}

func (model requestBodyModel) Update(msg tea.Msg) (requestBodyModel, tea.Cmd) {
	if !model.initialized {
		model.Init()
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		model.height = msg.Height

		return model, nil

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			o := model.input.Value()
			model.input.Reset()
			model.input.SetValue(o + encodedEnter)
			return model, nil
		}

	case requestBody:
		model.input.Reset()
		cleaned := strings.Replace(string(msg), "\n", encodedEnter, -1)
		model.input.SetValue(cleaned)

		return model, nil
	}

	var cmd tea.Cmd

	model.input, cmd = model.input.Update(msg)

	return model, cmd
}

func (model requestBodyModel) View() string {
	viewportStyle := lipgloss.NewStyle().Height(model.height).Padding(0, 2, 1, 2)
	decoded := strings.Replace(model.input.View(), encodedEnter, "\n", -1)

	return viewportStyle.Render(decoded)
}
