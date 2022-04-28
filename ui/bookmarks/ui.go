package bookmarks

import (
	"errors"
	"github.com/blackmann/gurl/lib"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"strings"
)

type Filter string

type Model struct {
	list        list.Model
	persistence lib.Persistence

	initialized bool
	bookmarks   []lib.Bookmark
}

func (m *Model) Init() tea.Cmd {
	listDefinition := list.New([]list.Item{}, lib.GetDefaultListDelegate(), 0, 0)
	listDefinition.SetShowTitle(false)
	listDefinition.SetFilteringEnabled(false)
	listDefinition.DisableQuitKeybindings()

	listDefinition.KeyMap.ShowFullHelp.Unbind()

	m.list = listDefinition
	m.initialized = true

	m.fetchBookmarks()

	return nil
}

func (m *Model) fetchBookmarks() {
	m.bookmarks = m.persistence.GetBookmarks()
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	if !m.initialized {
		m.Init()
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetHeight(msg.Height - 2)

	case Filter:
		var bookmarks []list.Item

		for _, b := range m.bookmarks {
			if strings.Contains(b.Name, string(msg)) {
				bookmarks = append(bookmarks, lib.ListItem{Key: b.Name, Value: b.Url, Ref: b})
			}
		}

		m.list.SetItems(bookmarks)

	case lib.Trigger:
		switch msg {
		case lib.UpdateBookmarks:
			m.fetchBookmarks()
			return m, nil
		}
	}

	m.list, _ = m.list.Update(msg)

	return m, nil
}

func (m Model) View() string {
	return lipgloss.NewStyle().Margin(1, 0, 1, 0).Render(m.list.View())
}

func (m Model) GetSelected() (lib.Bookmark, error) {
	if len(m.bookmarks) == 0 {
		return lib.Bookmark{}, errors.New("no bookmark entry")
	}

	return m.list.SelectedItem().(lib.ListItem).Ref.(lib.Bookmark), nil
}

func NewBookmarksList(persistence lib.Persistence) Model {
	return Model{persistence: persistence}
}
