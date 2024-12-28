package cli

import (
	"github.com/charmbracelet/bubbles/list"
)

func NewDefaultDelegate() list.ItemDelegate {
	d := list.NewDefaultDelegate()

	d.ShowDescription = true

	d.Styles.NormalTitle = d.Styles.NormalTitle.
		Foreground(fg0)
	d.Styles.NormalDesc = d.Styles.NormalDesc.
		Foreground(blue)

	d.Styles.SelectedTitle = d.Styles.SelectedTitle.
		Foreground(yellow).
		BorderForeground(yellow)

	d.Styles.SelectedDesc = d.Styles.SelectedDesc.
		Foreground(yellow).
		BorderForeground(yellow)

	return d
}
