package viewport

import (
	"github.com/blackmann/gurl/common"
	"github.com/blackmann/gurl/common/appcmd"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"strings"
)

type headerItem struct {
	key   string
	value string
}

func (h headerItem) FilterValue() string {
	return h.key
}

func (h headerItem) Title() string {
	return h.key
}

func (h headerItem) Description() string {
	return h.value
}

type responseHeadersModel struct {
	headersList list.Model
}

func NewResponseHeadersModel() responseHeadersModel {
	listDefinition := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	listDefinition.SetShowTitle(false)
	//listDefinition.DisableQuitKeybindings()

	listDefinition.KeyMap.ClearFilter.SetKeys("shift+esc")

	return responseHeadersModel{
		headersList: listDefinition,
	}
}

func (model responseHeadersModel) Update(msg tea.Msg) (responseHeadersModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		model.headersList.SetHeight(msg.Height)

	case common.Response:
		var items []list.Item
		for key, values := range msg.Headers {
			items = append(items, headerItem{key: key, value: strings.Join(values, ",")})
		}

		cmd := model.headersList.SetItems(items)

		return model, cmd

	case appcmd.Trigger:
		switch msg {
		case appcmd.LostFocus:
			model.headersList.SetFilteringEnabled(false)
			//model.headersList.SetShowFilter(false)
		case appcmd.GainFocus:
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
