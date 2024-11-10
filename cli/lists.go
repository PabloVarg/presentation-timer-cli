package cli

import (
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/lipgloss"
)

func NewDefaultList(d list.ItemDelegate) list.Model {
	l := list.New([]list.Item{}, d, 0, 0)

	l.Title = "Presentations"
	l.Styles.Title = lipgloss.NewStyle().Foreground(lipgloss.Color("#B8BB26"))

	l.SetSpinner(spinner.Dot)
	l.Styles.Spinner = l.Styles.Spinner.Foreground(lipgloss.Color("#00ee00"))

	l.SetShowTitle(true)

	l.StatusMessageLifetime = 5 * time.Second

	l.Styles.NoItems = lipgloss.NewStyle().Transform(func(s string) string { return "" })

	return l
}
