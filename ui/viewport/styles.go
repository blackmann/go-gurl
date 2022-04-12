package viewport

import "github.com/charmbracelet/lipgloss"

var (
	tabGroupStyle    = lipgloss.NewStyle().Margin(1, 1, 1, 2)
	inactiveTabStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#999")).
				MarginRight(2)
	activeTabStyle = lipgloss.NewStyle().MarginRight(2)
)
