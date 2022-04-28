package lib

import "github.com/charmbracelet/bubbletea"

func NavigateLeft() tea.Msg {
	return TabLeft
}

func NavigateRight() tea.Msg {
	return TabRight
}
