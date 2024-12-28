package cli

import "github.com/charmbracelet/lipgloss"

var (
	yellow = lipgloss.Color("#d8a657")
	orange = lipgloss.Color("#e78a4e")
	blue   = lipgloss.Color("#7daea3")
	fg0    = lipgloss.Color("#d4be98")
)

var (
	listContainerStyle = lipgloss.NewStyle().Padding(1)
	listTitleStyle     = lipgloss.NewStyle().Foreground(yellow)
	textInputFocused   = lipgloss.NewStyle()
)
