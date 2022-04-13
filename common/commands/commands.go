package commands

import tea "github.com/charmbracelet/bubbletea"

type FreeCommand string

type AppCommand int

var (
	NewRequest      AppCommand = 1
	CompleteRequest AppCommand = 2
)

func CreateFreeCommand(cmd string) tea.Cmd {
	return func() tea.Msg {
		return FreeCommand(cmd)
	}
}

func SubmitNewRequest() tea.Msg {
	return NewRequest
}
