package cli

import "github.com/charmbracelet/lipgloss"

var (
	Yellow      = lipgloss.Color("#d8a657")
	Green       = lipgloss.Color("#a9b665")
	Orange      = lipgloss.Color("#e78a4e")
	Blue        = lipgloss.Color("#7daea3")
	Red         = lipgloss.Color("#ea6962")
	Bg_diff_red = lipgloss.Color("#3c1f1e")
	Fg0         = lipgloss.Color("#d4be98")
	Fg1         = lipgloss.Color("#ddc7a1")
)

var (
	ContainerStyle         = lipgloss.NewStyle().Padding(1)
	CenteredContainerStyle = lipgloss.NewStyle().Align(lipgloss.Center).Padding(1)
	TitleStyle             = lipgloss.NewStyle().Foreground(Yellow)
	ErrorStyle             = lipgloss.NewStyle().Foreground(Red).Background(Bg_diff_red).Padding(1)
	// Lists
	ListTitleStyle = lipgloss.NewStyle().Foreground(Yellow)
	// Forms
	PromptStyle           = lipgloss.NewStyle().Foreground(Green)
	TextInputStyle        = lipgloss.NewStyle().Foreground(Fg0)
	TextInputFocusedStyle = lipgloss.NewStyle().Foreground(Fg1)
	CursorStyle           = lipgloss.NewStyle().Foreground(Fg0)
)
