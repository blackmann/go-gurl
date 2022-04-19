package viewport

import (
	"fmt"
	"github.com/blackmann/gurl/lib"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"strings"
)

type activeView int

var (
	INPUT activeView = 1
	LIST  activeView = 2
)

type headersModel struct {
	list        list.Model
	headerInput textinput.Model

	// states
	initialized      bool
	verticalPosition activeView
}

func (model *headersModel) Init() tea.Cmd {
	model.headerInput = textinput.New()

	model.headerInput.Placeholder = "Header-Key: Value"
	model.headerInput.Prompt = "+  "
	model.headerInput.Focus()

	listDefinition := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	listDefinition.SetShowTitle(false)
	model.list = listDefinition

	model.verticalPosition = INPUT
	model.initialized = true

	return nil
}

func (model headersModel) Update(msg tea.Msg) (headersModel, tea.Cmd) {
	if !model.initialized {
		model.Init()
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		model.list.SetHeight(msg.Height - 2)

	case requestHeaders:
		var items []list.Item
		for key, values := range msg {
			items = append(items, headerItem{key: key, value: strings.Join(values, ",")})
		}

		cmd := model.list.SetItems(items)

		return model, cmd

	case lib.Trigger:
		if model.verticalPosition == LIST {
			switch msg {
			case lib.LostFocus:
				model.list.SetFilteringEnabled(false)

			case lib.GainFocus:
				model.list.SetFilteringEnabled(true)
			}
		}

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			if model.verticalPosition == INPUT {
				parts := strings.Split(model.headerInput.Value(), ":")
				key, value := strings.Trim(parts[0], " "), strings.Trim(parts[1], " ")

				cmd := func() tea.Msg {
					return headerItem{key: key, value: value}
				}

				return model, cmd
			}

		case tea.KeyUp:
			if model.list.Cursor() == 0 {
				model.verticalPosition = INPUT
				return model, nil
			}

		case tea.KeyDown:
			if model.verticalPosition == INPUT && len(model.list.Items()) > 0 {
				model.verticalPosition = LIST

				return model, nil
			}
		}
	}

	var cmd tea.Cmd

	if model.verticalPosition == INPUT {
		model.headerInput, cmd = model.headerInput.Update(msg)
	} else {
		model.list, cmd = model.list.Update(msg)
	}

	return model, cmd
}

func (model headersModel) View() string {
	inputStyle := lipgloss.NewStyle().Margin(0, 2, 1, 2)
	input := inputStyle.Render(model.headerInput.View())
	return fmt.Sprintf("%s\n%s", input, model.list.View())
}
