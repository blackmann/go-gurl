package viewport

import (
	"github.com/blackmann/gurl/common/commands"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type keymap struct {
	nextTab     key.Binding
	previousTab key.Binding
}

type Model struct {
	activeTab int
	keybinds  keymap
	tabs      []string
	Height    int
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

func (model Model) View() string {
	viewportStyle := lipgloss.NewStyle().Height(model.Height)

	styledTabs := make([]string, len(model.tabs))

	for i, tab := range model.tabs {
		if i != model.activeTab {
			styledTabs = append(styledTabs, inactiveTabStyle.Render(tab))
		} else {
			styledTabs = append(styledTabs, activeTabStyle.Render(tab))
		}
	}

	tabsRow := tabGroupStyle.Render(lipgloss.JoinHorizontal(lipgloss.Left, styledTabs...))
	return viewportStyle.Render(tabsRow)
}

func (model Model) Update(msg tea.Msg) (Model, tea.Cmd) {
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

			return model, nil
		}

	case commands.FreeCommand:
		switch msg {
		case ":q":
			model.activeTab = 0

		case ":w":
			model.activeTab = 1

		case ":e":
			model.activeTab = 2
		}
	}

	return model, nil
}
