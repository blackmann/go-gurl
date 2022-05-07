package viewport

import (
	"fmt"
	"github.com/blackmann/go-gurl/lib"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"net/http"
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

	height int

	// state
	headers http.Header

	// tabs
	responseModel        responseModel
	requestBodyModel     requestBodyModel
	headersModel         headersModel
	responseHeadersModel responseHeadersModel
}

func NewViewport() Model {
	tabs := []string{"Headers", "Request Body", "Response", "Response Headers"}
	keybinds := keymap{
		nextTab:     key.NewBinding(key.WithKeys("shift+f"), key.WithHelp("⌥→", "Next tab")),
		previousTab: key.NewBinding(key.WithKeys("shift+b"), key.WithHelp("⌥←", "Prev tab")),
	}

	return Model{
		tabs:                 tabs,
		keybinds:             keybinds,
		responseHeadersModel: newResponseHeadersModel(),
		headers:              http.Header{},
	}
}

func (model *Model) SetResponse(response lib.Response) tea.Msg {
	return response
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

		default:
			var cmd tea.Cmd
			switch model.activeTab {
			case 0:
				model.headersModel, cmd = model.headersModel.Update(msg)
			case 1:
				model.requestBodyModel, cmd = model.requestBodyModel.Update(msg)
			case 2:
				model.responseModel, cmd = model.responseModel.Update(msg)
			case 3:
				model.responseHeadersModel, cmd = model.responseHeadersModel.Update(msg)
			}

			return model, cmd
		}

	case tea.WindowSizeMsg:
		renderHeight := msg.Height - 3
		resizeMsg := tea.WindowSizeMsg{Height: renderHeight, Width: msg.Width}

		model.responseModel, _ = model.responseModel.Update(resizeMsg)
		model.requestBodyModel, _ = model.requestBodyModel.Update(resizeMsg)
		model.headersModel, _ = model.headersModel.Update(resizeMsg)
		model.responseHeadersModel, _ = model.responseHeadersModel.Update(resizeMsg)

		return model, nil

	case lib.Pair:
		model.headers.Set(msg.Key, msg.Value)

		cmd := func() tea.Msg {
			return model.headers
		}

		return model, cmd

	case lib.Event:
		switch msg {
		case lib.TabLeft:
			if model.activeTab > 0 {
				model.activeTab -= 1
			}
			return model, nil

		case lib.TabRight:
			if model.activeTab < 3 {
				model.activeTab += 1
			}

			return model, nil
		}

	case lib.History:
		headers := http.Header{}
		for k, v := range msg.Headers {
			for _, value := range v {
				headers.Add(k, value)
			}
		}

		model.headers = headers

		model.headersModel, _ = model.headersModel.Update(headers)
		model.requestBodyModel, _ = model.requestBodyModel.Update(requestBody(msg.Body))
		model.responseModel.Reset()
		model.responseHeadersModel.Reset()
	}

	var cmds []tea.Cmd

	model.responseModel, cmd = model.responseModel.Update(msg)
	cmds = append(cmds, cmd)

	model.requestBodyModel, cmd = model.requestBodyModel.Update(msg)
	cmds = append(cmds, cmd)

	model.headersModel, cmd = model.headersModel.Update(msg)
	cmds = append(cmds, cmd)

	model.responseHeadersModel, cmd = model.responseHeadersModel.Update(msg)
	cmds = append(cmds, cmd)

	return model, tea.Batch(cmds...)
}

func (model Model) View() string {
	viewportStyle := lipgloss.NewStyle().Height(model.height)

	tabsCount := len(model.tabs)
	styledTabs := make([]string, tabsCount*2-1)

	for i, tab := range model.tabs {
		if i != model.activeTab {
			styledTabs = append(styledTabs, inactiveTabStyle.Render(tab))
		} else {
			styledTabs = append(styledTabs, activeTabStyle.Render(tab))
		}

		if i != len(model.tabs)-1 {
			styledTabs = append(styledTabs, inactiveTabStyle.Render(" | "))
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
	case 3:
		content = model.responseHeadersModel.View()
	}

	render := viewportStyle.Render(fmt.Sprintf("%s\n%s", tabsRow, content))

	return render
}

func (model Model) GetHeaders() http.Header {
	return model.headers
}

func (model *Model) GetBody() string {
	return model.requestBodyModel.Input()
}
