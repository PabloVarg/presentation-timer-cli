package cli

import (
	"github.com/charmbracelet/bubbles/list"
)

func NewDefaultDelegate() list.ItemDelegate {
	d := list.NewDefaultDelegate()

	d.ShowDescription = true

	d.Styles.NormalTitle = d.Styles.NormalTitle.
		Foreground(Fg0)
	d.Styles.NormalDesc = d.Styles.NormalDesc.
		Foreground(Blue)

	d.Styles.SelectedTitle = d.Styles.SelectedTitle.
		Foreground(Yellow).
		BorderForeground(Yellow)

	d.Styles.SelectedDesc = d.Styles.SelectedDesc.
		Foreground(Yellow).
		BorderForeground(Yellow)

	return d
}
