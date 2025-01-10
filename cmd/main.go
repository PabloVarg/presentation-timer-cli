package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/PabloVarg/presentation-timer-cli/cli"
	"github.com/PabloVarg/presentation-timer-cli/cli/presentations"
	"github.com/PabloVarg/presentation-timer-cli/internal/api"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	run()
}

func run() {
	readEnv()

	pm := cli.ProgramModel{
		Logger: getLogger(),
	}

	p := tea.NewProgram(
		presentations.NewListPresentations(pm),
		tea.WithAltScreen(),
	)
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

func getLogger() *slog.Logger {
	f, err := os.OpenFile(
		"./cli.log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0664,
	)
	if err != nil {
		panic("could not open log file")
	}

	return slog.New(slog.NewJSONHandler(f, nil))
}
