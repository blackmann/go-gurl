package history

import (
	"fmt"
	"github.com/blackmann/gurl/lib"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/dustin/go-humanize"
)

type Filter string

type Model struct {
	persistence lib.Persistence
	list        list.Model

	history []lib.History

	initialized bool
}

func (m *Model) OnChange(persistence lib.Persistence) {
	//TODO implement me
	panic("implement me")
}

func (m *Model) Init() tea.Cmd {
	delegate := lib.GetDefaultListDelegate()

	listDefinition := list.New([]list.Item{}, delegate, 0, 0)
	listDefinition.SetShowTitle(false)
	listDefinition.SetFilteringEnabled(false)
	listDefinition.DisableQuitKeybindings()

	listDefinition.KeyMap.ShowFullHelp.Unbind()

	m.list = listDefinition
	m.initialized = true

	m.persistence.AddListener(m)

	m.prefetchHistory()

	return nil
}

func (m *Model) prefetchHistory() {
	m.history = m.persistence.GetHistory()
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	if !m.initialized {
		m.Init()
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetHeight(msg.Height - 2)

	case Filter:
		var historyItems []list.Item
		for _, item := range m.history {
			historyItems = append(historyItems,
				lib.ListItem{
					Key:   fmt.Sprintf("%s %s", item.Method, item.Url),
					Value: fmt.Sprintf("%s $%d", humanize.Time(item.Date), item.ID),
				})
		}

		m.list.SetItems(historyItems)
	}

	m.list, _ = m.list.Update(msg)

	return m, nil
}

func (m Model) View() string {
	return lipgloss.NewStyle().Margin(1, 0, 1, 0).Render(m.list.View())
}

func NewHistory(persistence lib.Persistence) Model {
	h := Model{persistence: persistence}

	return h
}