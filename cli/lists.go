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
	l.SetShowTitle(true)
	l.Styles.Title = listTitleStyle

	l.SetSpinner(spinner.Dot)
	l.Styles.Spinner = l.Styles.Spinner.Foreground(lipgloss.Color("#00ee00"))

	l.StatusMessageLifetime = 5 * time.Second

	l.Styles.NoItems = lipgloss.NewStyle().Transform(func(s string) string { return "" })

	return l
}
