package appcmd

import tea "github.com/charmbracelet/bubbletea"

type FreeText string

type Trigger int

type Response struct {
	Body string
}

var (
	NewRequest Trigger = 1
)

func SubmitNewRequest() tea.Msg {
	return NewRequest
}
