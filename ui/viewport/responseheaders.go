package viewport

import (
	"sort"
	"strings"

	"github.com/blackmann/go-gurl/lib"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type responseHeadersModel struct {
	headersList list.Model
}

func newResponseHeadersModel() tea.Model {
	listDefinition := list.New([]list.Item{}, lib.GetDefaultListDelegate(), 0, 0)
	listDefinition.SetShowTitle(false)
	listDefinition.SetFilteringEnabled(false)
	listDefinition.DisableQuitKeybindings()

	listDefinition.KeyMap.ShowFullHelp.Unbind()

	return &responseHeadersModel{
		headersList: listDefinition,
	}
}

func (model responseHeadersModel) Init() tea.Cmd {
	return nil
}

func (model responseHeadersModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		model.headersList.SetSize(msg.Width-2, msg.Height)

	case lib.Response:
		var items []list.Item
		for key, values := range msg.Headers {
			items = append(items, lib.Pair{
				Key:   key,
				Value: strings.Join(values, ","),
			})
		}

		sort.Slice(items, func(i, j int) bool {
			var (
				headerA, okA = items[i].(lib.Pair)
				headerB, okB = items[j].(lib.Pair)
			)

			if !okA && !okB {
				return false
			} else if !okA {
				return strings.Compare("", headerB.Title()) < 0
			} else if !okB {
				return strings.Compare(headerA.Title(), "") < 0
			}

			return strings.Compare(headerA.Title(), headerB.Title()) < 0
		})

		cmd := model.headersList.SetItems(items)

		return model, cmd

	case lib.Event:
		if msg == lib.Reset {
			model.headersList.SetItems([]list.Item{})

			return model, nil
		}
	}

	var cmd tea.Cmd
	model.headersList, cmd = model.headersList.Update(msg)

	return model, cmd
}

func (model responseHeadersModel) View() string {
	return model.headersList.View()
}
