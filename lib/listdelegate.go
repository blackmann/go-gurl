package lib

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
)

func GetDefaultListDelegate() list.DefaultDelegate {
	delegate := list.NewDefaultDelegate()
	delegate.Styles.SelectedTitle = delegate.Styles.SelectedTitle.Copy().
		BorderForeground(lipgloss.Color(ANSIRed)).
		Foreground(lipgloss.Color(ANSIRed))

	delegate.Styles.SelectedDesc = delegate.Styles.SelectedDesc.Copy().
		BorderForeground(lipgloss.Color(ANSIRed)).
		Foreground(lipgloss.Color(ANSIMagenta))

	return delegate
}
