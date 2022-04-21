package statusbar

import (
	"github.com/blackmann/gurl/lib"
	"github.com/charmbracelet/lipgloss"
)

var (
	barStyle = lipgloss.NewStyle().
			Padding(0, 1, 0, 1).
			MarginTop(1)

	neutralStatusStyle = lipgloss.NewStyle().Background(lipgloss.Color("#fff")).
				Foreground(lipgloss.Color("#000")).
				Padding(0, 1, 0, 1)

	errorStatusStyle = lipgloss.NewStyle().Background(lipgloss.Color("#fff")).
				Background(lib.ColorDanger).
				Foreground(lipgloss.Color("#fff")).
				Padding(0, 1, 0, 1)

	okStatusStyle = lipgloss.NewStyle().Background(lipgloss.Color("#fff")).
			Background(lib.ColorOk).
			Foreground(lipgloss.Color("#000")).
			Padding(0, 1, 0, 1)
)
