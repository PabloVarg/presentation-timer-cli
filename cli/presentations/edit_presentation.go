package presentations

import (
	"os"
	"strings"

	"github.com/PabloVarg/presentation-timer-cli/cli"
	"github.com/PabloVarg/presentation-timer-cli/internal/api"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type EditPresentation struct {
	cli.ProgramModel
	cli.APIModel
	cli.StyledComponent
	cli.FormModel

	isLoading      bool
	PresentationID int
}

var editPresentationInputs = map[string]int{
	"name": 0,
}

func NewEditPresentation(m cli.ProgramModel, ID int) EditPresentation {
	nameInput := cli.NewDefaultTextInput()
	nameInput.Placeholder = "My Presentation"
	nameInput.Prompt = "Name: "

	nameInput.Focus()

	return EditPresentation{
		PresentationID: ID,
		ProgramModel:   m,
		FormModel: cli.FormModel{
			Inputs: []textinput.Model{
				nameInput,
			},
		},
	}
}

func (m EditPresentation) Init() tea.Cmd {
	m.isLoading = true

	return tea.Batch(
		textinput.Blink,
		func() tea.Msg {
			p, err := api.GetPresentation(m.Api, os.LookupEnv, m.PresentationID)
			if err != nil {
				return cli.FetchError{
					Err: err.Error(),
				}
			}

			return p
		},
	)
}

func (m EditPresentation) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := make([]tea.Cmd, 0)
	cmds = append(cmds, m.ProgramModel.Update(msg))

	switch msg := msg.(type) {
	case api.Presentation:
		m.Inputs[editPresentationInputs["name"]].SetValue(msg.Name)
		m.isLoading = false
	case cli.FetchError:
		m.Inputs = nil
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit
		case tea.KeyEsc:
			return cli.Transition(NewListPresentations(m.ProgramModel))
		case tea.KeyEnter:
			cmds = append(cmds, func() tea.Msg {
				err := api.UpdatePresentation(m.Api, os.LookupEnv, api.EditPresentationMsg{
					ID:   m.PresentationID,
					Name: m.Inputs[editPresentationInputs["name"]].Value(),
				})
				if err != nil {
					return cli.FormError{
						Err: err.Error(),
					}
				}
				return tea.KeyMsg{
					Type: tea.KeyEsc,
				}
			})
		}
	}

	m.UpdateForm(msg, tea.KeyEnter)

	cmds = append(cmds, m.UpdateInputs(msg))
	return m, tea.Batch(cmds...)
}

func (m EditPresentation) View() string {
	var sb strings.Builder

	sb.WriteString(
		cli.CenteredContainerStyle.Width(m.Width).
			Render(cli.TitleStyle.Render("Edit Presentation")),
	)
	sb.WriteRune('\n')

	for i := range m.Inputs {
		sb.WriteString(m.Inputs[i].View())
		sb.WriteRune('\n')
	}

	if m.Err != nil {
		sb.WriteRune('\n')
		sb.WriteString(cli.ErrorStyle.Render(m.Err.Error()))
	}

	return cli.ContainerStyle.Render(sb.String())
}
