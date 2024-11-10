package cli

import (
	"fmt"
	"log"

	"github.com/PabloVarg/presentation-timer-cli/internal/api"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ListPresentations struct {
	Spinner spinner.Model
	APIModel
	APIResponseModel[[]api.Presentation]
}

func NewListPresentations() ListPresentations {
	s := spinner.New()
	s.Spinner = spinner.Points

	return ListPresentations{
		Spinner: s,
	}
}

func (m ListPresentations) Init() tea.Cmd {
	return tea.Batch(func() tea.Msg {
		result, err := api.GetPresentations(m.api, os.LookupEnv)
		if err != nil {
			return err
		}

		return result
	}, m.Spinner.Tick)
}

func (m ListPresentations) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case error:
		m.SetErr(msg)
	case []api.Presentation:
		m.SetData(msg)
	case tea.KeyMsg:
		m.handleKeyPress(msg)
	case spinner.TickMsg:
		if m.done {
			return m, nil
		}

		s, cmd := m.Spinner.Update(msg)
		m.Spinner = s
		return m, cmd
	}

	return m, nil
}

func (m ListPresentations) View() string {
	resStyle := lipgloss.
		NewStyle().
		Foreground(lipgloss.Color("#00ee00"))
	errStyle := lipgloss.
		NewStyle().
		Foreground(lipgloss.Color("#ee0000"))

	if m.Loading() {
		return m.Spinner.View()
	}

	if m.err != nil {
		return errStyle.Render(m.err.Error())
	}

	return resStyle.Render(fmt.Sprintf("%+v", m.data))
}

func (m ListPresentations) handleKeyPress(key tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch key.String() {
	case "j":
		log.Println("down")
	case "k":
		log.Println("up")
	}

	return m, nil
}
