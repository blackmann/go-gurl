package viewport

import "github.com/charmbracelet/lipgloss"

var (
	activeTabStyle = lipgloss.NewStyle().MarginRight(2)

	contentStyle = lipgloss.NewStyle().Margin(0, 2, 0, 1)

	disabledContent = lipgloss.NewStyle().Foreground(lipgloss.Color("#999"))

	inactiveTabStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#999")).
				MarginRight(2)

	tabGroupStyle = lipgloss.NewStyle().Margin(1, 1, 1, 2)
)
