package commands

import tea "github.com/charmbracelet/bubbletea"

type FreeCommand string

func CreateFreeCommand(cmd string) tea.Cmd {
	return func() tea.Msg {
		return FreeCommand(cmd)
	}
}
