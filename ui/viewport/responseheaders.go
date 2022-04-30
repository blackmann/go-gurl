package viewport

import (
	"github.com/blackmann/go-gurl/lib"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"strings"
)

type responseHeadersModel struct {
	headersList list.Model
}

func newResponseHeadersModel() responseHeadersModel {
	listDefinition := list.New([]list.Item{}, lib.GetDefaultListDelegate(), 0, 0)
	listDefinition.SetShowTitle(false)
	listDefinition.SetFilteringEnabled(false)
	listDefinition.DisableQuitKeybindings()

	listDefinition.KeyMap.ShowFullHelp.Unbind()

	return responseHeadersModel{
		headersList: listDefinition,
	}
}

func (model responseHeadersModel) Update(msg tea.Msg) (responseHeadersModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		model.headersList.SetSize(msg.Width-2, msg.Height)

	case lib.Response:
		var items []list.Item
		for key, values := range msg.Headers {
			items = append(items, lib.ListItem{
				Key:   key,
				Value: strings.Join(values, ","),
			})
		}

		cmd := model.headersList.SetItems(items)

		return model, cmd
	}

	var cmd tea.Cmd
	model.headersList, cmd = model.headersList.Update(msg)

	return model, cmd
}

func (model responseHeadersModel) View() string {
	return model.headersList.View()
}

func (model *responseHeadersModel) Reset() {
	model.headersList.SetItems([]list.Item{})
}
