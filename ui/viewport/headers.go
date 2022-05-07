package viewport

import (
	"fmt"
	"github.com/blackmann/go-gurl/lib"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"net/http"
	"sort"
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

func (model *headersModel) initialize() {
	model.headerInput = textinput.New()

	model.headerInput.Placeholder = "Header-Key: Value"
	model.headerInput.Prompt = "+  "
	model.headerInput.Focus()

	listDefinition := list.New([]list.Item{}, lib.GetDefaultListDelegate(), 0, 0)
	listDefinition.SetShowTitle(false)
	listDefinition.SetFilteringEnabled(false)
	listDefinition.DisableQuitKeybindings()

	listDefinition.KeyMap.ShowFullHelp.Unbind()

	model.list = listDefinition

	model.verticalPosition = INPUT
	model.initialized = true
}

func (model headersModel) Update(msg tea.Msg) (headersModel, tea.Cmd) {
	if !model.initialized {
		model.initialize()
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		model.list.SetHeight(msg.Height - 2)

	case http.Header:
		var items []list.Item
		for key, values := range msg {
			items = append(items, lib.Pair{Key: key, Value: strings.Join(values, ",")})
		}

		// Sort based on entry
		sort.Slice(items, func(i, j int) bool {
			return j < i
		})

		cmd := model.list.SetItems(items)

		return model, cmd

	case lib.Event:
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
					return lib.Pair{Key: key, Value: value}
				}

				model.headerInput.SetValue("")

				return model, cmd
			}

		case tea.KeyUp:
			if model.list.Cursor() == 0 {
				model.verticalPosition = INPUT
				return model, textinput.Blink
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

	if model.verticalPosition == LIST {
		inputStyle = inputStyle.Foreground(lipgloss.Color("#999"))
	}

	input := inputStyle.Render(model.headerInput.View())
	return fmt.Sprintf("%s\n%s", input, model.list.View())
}
