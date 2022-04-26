package bookmarks

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Filter string

type Model struct {
	list list.Model

	initialized bool
}

func (m *Model) Init() tea.Cmd {
	listDefinition := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	listDefinition.SetShowTitle(false)
	listDefinition.SetFilteringEnabled(false)
	listDefinition.DisableQuitKeybindings()

	listDefinition.KeyMap.ShowFullHelp.Unbind()

	m.list = listDefinition
	m.initialized = true

	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	if !m.initialized {
		m.Init()
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetHeight(msg.Height - 2)

	case Filter:
		// TODO:
	}

	m.list, _ = m.list.Update(msg)

	return m, nil
}

func (m Model) View() string {
	return lipgloss.NewStyle().Margin(1, 0, 1, 0).Render(m.list.View())
}
