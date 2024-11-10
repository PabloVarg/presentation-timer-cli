package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/PabloVarg/presentation-timer-cli/cli"
	"github.com/PabloVarg/presentation-timer-cli/internal/api"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	run()
}

func run() {
	readEnv()

	p := tea.NewProgram(cli.NewListPresentations(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprint(os.Stderr, err)
	}
}

func readEnv() {
	url := flag.String(
		"url",
		"http://localhost:8000",
		"URL of backend",
	)
	os.Setenv(api.API_URL_KEY, *url)
}
