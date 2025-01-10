package cli

import (
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/lipgloss"
)

func NewDefaultList(d list.ItemDelegate, helpKeys []key.Binding, title string) list.Model {
	l := list.New([]list.Item{}, d, 0, 0)

	l.Title = title
	l.SetShowTitle(true)
	l.Styles.Title = ListTitleStyle

	l.SetSpinner(spinner.Dot)
	l.Styles.Spinner = l.Styles.Spinner.Foreground(lipgloss.Color("#00ee00"))

	l.Styles.NoItems = lipgloss.NewStyle().Transform(func(s string) string { return "" })

	l.StatusMessageLifetime = 5 * time.Second
	l.AdditionalFullHelpKeys = func() []key.Binding {
		return helpKeys
	}

	return l
}
