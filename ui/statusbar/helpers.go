package statusbar

import (
	"github.com/blackmann/gurl/common/status"
	tea "github.com/charmbracelet/bubbletea"
)

func CommandMsg(command string) tea.Msg {
	return commandInput(command)
}

func StatusMsg(status status.Status) tea.Msg {
	return statusUpdate(status)
}

func Width(width int) tea.Msg {
	return widthUpdate(width)
}