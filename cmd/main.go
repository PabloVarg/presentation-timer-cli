package main

import (
	"flag"
	"fmt"
	"log"
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

	f, err := os.OpenFile("log.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		panic("could not open log file")
	}
	log.SetOutput(f)
	log.SetFlags(log.Ldate | log.Ltime)
}
