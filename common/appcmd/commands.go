package appcmd

import tea "github.com/charmbracelet/bubbletea"

type FreeText string

type Trigger int

var (
	NewRequest Trigger = 1
	LostFocus  Trigger = 2
	GainFocus  Trigger = 3
)

func SubmitNewRequest() tea.Msg {
	return NewRequest
}
