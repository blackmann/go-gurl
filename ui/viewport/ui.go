package viewport

import (
	"fmt"
	"github.com/blackmann/gurl/common"
	"github.com/blackmann/gurl/common/appcmd"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type keymap struct {
	nextTab     key.Binding
	previousTab key.Binding
}

type Model struct {
	activeTab    int
	keybinds     keymap
	tabs         []string
	responseBody string
	enabled      bool

	height int
	Width  int

	// tabs
	responseModel
	requestBodyModel
	headersModel
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

func (model *Model) SetResponse(response common.Response) tea.Msg {
	return response
}

func (model *Model) SetEnabled(enabled bool) {
	model.enabled = enabled
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

		return model, nil

	case tea.WindowSizeMsg:
		renderHeight := msg.Height - 3
		resizeMsg := tea.WindowSizeMsg{Height: renderHeight, Width: msg.Width}

		model.responseModel, _ = model.responseModel.Update(resizeMsg)
		model.requestBodyModel, _ = model.requestBodyModel.Update(resizeMsg)
		model.headersModel, _ = model.headersModel.Update(resizeMsg)

		return model, nil
	}

	var cmds []tea.Cmd

	model.responseModel, cmd = model.responseModel.Update(msg)
	cmds = append(cmds, cmd)

	model.requestBodyModel, cmd = model.requestBodyModel.Update(msg)
	cmds = append(cmds, cmd)

	model.headersModel, cmd = model.headersModel.Update(msg)
	cmds = append(cmds, cmd)

	return model, tea.Batch(cmds...)
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

	content := ""

	switch model.activeTab {
	case 0:
		content = model.headersModel.View()
	case 1:
		content = model.requestBodyModel.View()
	case 2:
		content = model.responseModel.View()
	}

	return viewportStyle.Render(fmt.Sprintf("%s\n%s", tabsRow, content))
}
