package main

import (
	"fmt"
	"os"

	"github.com/PabloVarg/presentation-timer-cli/cli"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	run()
}

func run() {
	p := tea.NewProgram(cli.NewListPresentations(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprint(os.Stderr, err)
	}
}
