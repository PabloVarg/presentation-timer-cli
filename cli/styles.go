package cli

import "github.com/charmbracelet/lipgloss"

var (
	yellow = lipgloss.Color("#d8a657")
	orange = lipgloss.Color("#e78a4e")
	blue   = lipgloss.Color("#7daea3")
	red    = lipgloss.Color("#ea6962")
	fg0    = lipgloss.Color("#d4be98")
)

var (
	containerStyle = lipgloss.NewStyle().Padding(1)
	titleStyle     = lipgloss.NewStyle().Foreground(yellow)
	errorStyle     = lipgloss.NewStyle().Foreground(red)
	// Lists
	listTitleStyle = lipgloss.NewStyle().Foreground(yellow)
	// Forms
	textInputFocused = lipgloss.NewStyle().Foreground(fg0)
)
