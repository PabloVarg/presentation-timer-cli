package cli

import (
	"fmt"
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
	nameInput := textinput.New()
	nameInput.Placeholder = "Name"

	nameInput.Focus()

	return CreatePresentation{
		ProgramModel: m,
		FormModel: FormModel{
			inputs: []textinput.Model{
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
		case tea.KeyCtrlC, tea.KeyEsc:
			nextModel := NewListPresentations(m.ProgramModel)
			return nextModel, nextModel.Init()
		case tea.KeyEnter:
			cmds = append(cmds, func() tea.Msg {
				err := api.CreatePresentation(m.api, os.LookupEnv, api.CreatePresentationMsg{
					Name: m.inputs[0].Value(),
				})
				if err != nil {
					return fmt.Errorf("error creating presentation: %s", err)
				}

				return nil
			})
		}
	}

	cmds = append(cmds, m.updateInputs(msg))
	return m, tea.Batch(cmds...)
}

func (m CreatePresentation) View() string {
	var sb strings.Builder

	sb.WriteString(titleStyle.Render("Create a Presentation"))
	sb.WriteRune('\n')

	for i := range m.inputs {
		sb.WriteString(m.inputs[i].View())
		sb.WriteRune('\n')
	}

	return sb.String()
}
