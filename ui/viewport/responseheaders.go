package viewport

import (
	"github.com/blackmann/gurl/lib"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"strings"
)

type responseHeadersModel struct {
	headersList list.Model
}

func newResponseHeadersModel() responseHeadersModel {
	listDefinition := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	listDefinition.SetShowTitle(false)

	return responseHeadersModel{
		headersList: listDefinition,
	}
}

func (model responseHeadersModel) Update(msg tea.Msg) (responseHeadersModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		model.headersList.SetHeight(msg.Height)

	case lib.Response:
		var items []list.Item
		for key, values := range msg.Headers {
			items = append(items, headerItem{key: key, value: strings.Join(values, ",")})
		}

		cmd := model.headersList.SetItems(items)

		return model, cmd

	case lib.Trigger:
		switch msg {
		case lib.LostFocus:
			model.headersList.SetFilteringEnabled(false)
			//model.headersList.SetShowFilter(false)
		case lib.GainFocus:
			model.headersList.SetFilteringEnabled(true)
			model.headersList.SetShowFilter(true)
		}

	}

	var cmd tea.Cmd
	model.headersList, cmd = model.headersList.Update(msg)

	return model, cmd
}

func (model responseHeadersModel) View() string {
	return model.headersList.View()
}
