package lib

import "github.com/charmbracelet/bubbletea"

func SubmitNewRequest() tea.Msg {
	return NewRequest
}

func NavigateLeft() tea.Msg {
	return TabLeft
}

func NavigateRight() tea.Msg {
	return TabRight
}
