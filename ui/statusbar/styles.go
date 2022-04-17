package statusbar

import "github.com/charmbracelet/lipgloss"

var (
	barStyle = lipgloss.NewStyle().
			Padding(0, 1, 0, 1)

	idleStatusStyle = lipgloss.NewStyle().Background(lipgloss.Color("#fff")).
			Foreground(lipgloss.Color("#000")).
			Padding(0, 1, 0, 1)
)
