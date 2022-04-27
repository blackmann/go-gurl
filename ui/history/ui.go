package history

import (
	"fmt"
	"github.com/blackmann/gurl/lib"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/dustin/go-humanize"
	"gorm.io/gorm/utils"
	"strings"
)

type Filter string

type Model struct {
	persistence lib.Persistence
	list        list.Model

	history []lib.History

	initialized bool
	filter      Filter
	width       int
}

func (m *Model) Init() tea.Cmd {
	delegate := lib.GetDefaultListDelegate()

	listDefinition := list.New([]list.Item{}, delegate, 0, 0)
	listDefinition.SetShowTitle(false)
	listDefinition.SetFilteringEnabled(false)
	listDefinition.DisableQuitKeybindings()

	listDefinition.KeyMap.ShowFullHelp.Unbind()

	m.list = listDefinition

	m.prefetchHistory()
	m.initialized = true

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
		m.list.SetSize(msg.Width-2, msg.Height-2)
		m.width = msg.Width

		return m, nil

	case Filter:
		var historyItems []list.Item

		for _, item := range m.history {
			idString := utils.ToString(item.ID)
			if !(strings.HasPrefix(idString, string(msg)) ||
				strings.Contains(item.Annotation, string(msg))) {
				continue
			}

			left := lipgloss.NewStyle().
				Width(m.width/2 - 4).
				Render(fmt.Sprintf("%d", item.Status))

			right := lipgloss.NewStyle().
				Width(m.width/2 - 4).
				Align(lipgloss.Right).
				Render(fmt.Sprintf("%s $%d", humanize.Time(item.Date), item.ID))

			historyItems = append(historyItems,
				lib.ListItem{
					Key:   fmt.Sprintf("%s %s", item.Method, item.Url),
					Value: fmt.Sprintf("%s%s", left, right),
				})
		}

		m.list.SetItems(historyItems)
		return m, nil

	case lib.Trigger:
		if msg == lib.NewHistory {
			m.history = m.persistence.GetHistory()
			m.list.Select(0)
			return m, nil
		}
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
