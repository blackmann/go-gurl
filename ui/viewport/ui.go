package viewport

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/blackmann/gurl/common/appcmd"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type keymap struct {
	nextTab     key.Binding
	previousTab key.Binding
}

type Model struct {
	viewport     viewport.Model
	activeTab    int
	keybinds     keymap
	tabs         []string
	responseBody string
	enabled      bool

	height int
	Width  int
}

func NewViewport() Model {
	tabs := []string{"Headers (:q)", "Request Body (:w)", "Response (:e)"}
	keybinds := keymap{
		nextTab:     key.NewBinding(key.WithKeys("shift+f"), key.WithHelp("⌥→", "Next tab")),
		previousTab: key.NewBinding(key.WithKeys("shift+b"), key.WithHelp("⌥←", "Prev tab")),
	}

	return Model{
		tabs:     tabs,
		keybinds: keybinds,
	}
}

func (model *Model) SetResponse(response string) {
	model.responseBody = response
}

func (model *Model) SetEnabled(enabled bool) {
	model.enabled = enabled
}

func (model *Model) SetHeight(height int) {
	model.height = height
	model.viewport.Height = height - 4
}

func (model Model) getContent() string {
	content := ""

	switch model.activeTab {
	case 2:
		var pretty bytes.Buffer
		if err := json.Indent(&pretty, []byte(model.responseBody), "", "  "); err == nil {
			content = pretty.String()
		} else {
			content = model.responseBody
		}

		content = lipgloss.NewStyle().Width(model.Width-2).
			Margin(0, 1, 0, 1).
			Render(content)

		if !model.enabled {
			content = disabledContent.Render(content)
		}
	}

	return content
}

func (model Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, model.keybinds.nextTab):
			model.activeTab = (model.activeTab + 1) % len(model.tabs)
			return model, nil

		case key.Matches(msg, model.keybinds.previousTab):
			var newTab int
			if model.activeTab == 0 {
				newTab = 2
			} else {
				newTab = model.activeTab - 1
			}
			model.activeTab = newTab

		default:
			model.viewport, cmd = model.viewport.Update(msg)

			return model, cmd
		}

	case appcmd.FreeText:
		switch msg {
		case ":q":
			model.activeTab = 0

		case ":w":
			model.activeTab = 1

		case ":e":
			model.activeTab = 2
		}
	}

	content := model.getContent()
	model.viewport.SetContent(content)

	return model, nil
}

func (model Model) View() string {
	viewportStyle := lipgloss.NewStyle().Height(model.height)

	styledTabs := make([]string, len(model.tabs))

	for i, tab := range model.tabs {
		if i != model.activeTab {
			styledTabs = append(styledTabs, inactiveTabStyle.Render(tab))
		} else {
			styledTabs = append(styledTabs, activeTabStyle.Render(tab))
		}
	}

	tabsRow := tabGroupStyle.Render(lipgloss.JoinHorizontal(lipgloss.Left, styledTabs...))

	return viewportStyle.Render(fmt.Sprintf("%s\n%s", tabsRow, model.viewport.View()))
}
