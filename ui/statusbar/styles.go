package statusbar

import (
	"github.com/blackmann/go-gurl/lib"
	"github.com/charmbracelet/lipgloss"
)

var (
	barStyle = lipgloss.NewStyle().
			Padding(0, 1, 0, 1).
			MarginTop(1)

	neutralStatusStyle = lipgloss.NewStyle().
				Background(lipgloss.Color(lib.ANSIWhite)).
				Foreground(lipgloss.Color(lib.ANSIBlack)).
				Padding(0, 1, 0, 1)

	errorStatusStyle = lipgloss.NewStyle().
				Background(lipgloss.Color(lib.ANSIRed)).
				Foreground(lipgloss.Color("#fff")).
				Padding(0, 1, 0, 1)

	okStatusStyle = lipgloss.NewStyle().
			Background(lipgloss.Color(lib.ANSIGreen)).
			Foreground(lipgloss.Color(lib.ANSIBlack)).
			Padding(0, 1, 0, 1)
)
