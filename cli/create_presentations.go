package cli

import (
	"os"
	"strings"

	"github.com/PabloVarg/presentation-timer-cli/internal/api"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type CreatePresentation struct {
	ProgramModel
	APIModel
	StyledComponent
	FormModel
}

func NewCreatePresentation(m ProgramModel) CreatePresentation {
	nameInput := NewDefaultTextInput()
	nameInput.Placeholder = "My Presentation"
	nameInput.Prompt = "Name: "

	nameInput.Focus()

	return CreatePresentation{
		ProgramModel: m,
		FormModel: FormModel{
			Inputs: []textinput.Model{
				nameInput,
			},
		},
	}
}

func (m CreatePresentation) Init() tea.Cmd {
	return textinput.Blink
}

func (m CreatePresentation) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := make([]tea.Cmd, 0)
	cmds = append(cmds, m.ProgramModel.Update(msg))

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit
		case tea.KeyEsc:
			return Transition(NewListPresentations(m.ProgramModel))
		case tea.KeyEnter:
			cmds = append(cmds, func() tea.Msg {
				err := api.CreatePresentation(m.Api, os.LookupEnv, api.CreatePresentationMsg{
					Name: m.Inputs[0].Value(),
				})
				if err != nil {
					return FormError{
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

func (m CreatePresentation) View() string {
	var sb strings.Builder

	sb.WriteString(
		CenteredContainerStyle.Width(m.Width).
			Render(TitleStyle.Render("Create a Presentation")),
	)
	sb.WriteRune('\n')

	for i := range m.Inputs {
		sb.WriteString(m.Inputs[i].View())
		sb.WriteRune('\n')
	}

	if m.Err != nil {
		sb.WriteRune('\n')
		sb.WriteString(ErrorStyle.Render(m.Err.Error()))
	}

	return ContainerStyle.Render(sb.String())
}
