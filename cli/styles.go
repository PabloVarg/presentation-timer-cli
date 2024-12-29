package cli

import "github.com/charmbracelet/lipgloss"

var (
	yellow      = lipgloss.Color("#d8a657")
	green       = lipgloss.Color("#a9b665")
	orange      = lipgloss.Color("#e78a4e")
	blue        = lipgloss.Color("#7daea3")
	red         = lipgloss.Color("#ea6962")
	bg_diff_red = lipgloss.Color("#3c1f1e")
	fg0         = lipgloss.Color("#d4be98")
	fg1         = lipgloss.Color("#ddc7a1")
)

var (
	containerStyle         = lipgloss.NewStyle().Padding(1)
	centeredContainerStyle = lipgloss.NewStyle().Align(lipgloss.Center).Padding(1)
	titleStyle             = lipgloss.NewStyle().Foreground(yellow)
	errorStyle             = lipgloss.NewStyle().Foreground(red).Background(bg_diff_red).Padding(1)
	// Lists
	listTitleStyle = lipgloss.NewStyle().Foreground(yellow)
	// Forms
	promptStyle           = lipgloss.NewStyle().Foreground(green)
	textInputStyle        = lipgloss.NewStyle().Foreground(fg0)
	textInputFocusedStyle = lipgloss.NewStyle().Foreground(fg1)
	cursorStyle           = lipgloss.NewStyle().Foreground(fg0)
)
