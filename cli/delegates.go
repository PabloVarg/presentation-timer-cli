package cli

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
)

func NewDefaultDelegate() list.ItemDelegate {
	d := list.NewDefaultDelegate()

	d.ShowDescription = true

	d.Styles.SelectedTitle = d.Styles.SelectedTitle.
		Foreground(lipgloss.Color("#FABD2F")).
		BorderForeground(lipgloss.Color("#D79921"))

	d.Styles.SelectedDesc = d.Styles.SelectedDesc.
		Foreground(lipgloss.Color("#D79921")).
		BorderForeground(lipgloss.Color("#D79921"))

	return d
}
