package cli

import (
	"os"
	"strings"
	"time"

	"github.com/PabloVarg/presentation-timer-cli/internal/api"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type CreateSection struct {
	ProgramModel
	APIModel
	StyledComponent
	FormModel

	presentationID int
}

var createSectionInputs = map[string]int{
	"name":     0,
	"duration": 1,
}

func NewCreateSection(m ProgramModel, presentationID int) CreateSection {
	nameInput := NewDefaultTextInput()
	nameInput.Placeholder = "My Presentation"
	nameInput.Prompt = "Name: "

	durationInput := NewDefaultTextInput()
	durationInput.Placeholder = "15m20s"
	durationInput.Prompt = "Duration: "

	nameInput.Focus()

	return CreateSection{
		presentationID: presentationID,
		ProgramModel:   m,
		FormModel: FormModel{
			Inputs: []textinput.Model{
				nameInput,
				durationInput,
			},
		},
	}
}

func (m CreateSection) Init() tea.Cmd {
	return textinput.Blink
}

func (m CreateSection) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := make([]tea.Cmd, 0)
	cmds = append(cmds, m.ProgramModel.Update(msg))

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit
		case tea.KeyEsc:
			return Transition(NewListSections(m.ProgramModel, m.presentationID))
		case tea.KeyEnter:
			cmds = append(cmds, func() tea.Msg {
				d, err := time.ParseDuration(m.Inputs[createSectionInputs["duration"]].Value())
				if err != nil {
					return FormError{
						Err: "Wrong format for duration",
					}
				}

				err = api.CreateSection(m.Api, os.LookupEnv, api.CreateSectionMsg{
					PresentationID: m.presentationID,
					Name:           m.Inputs[createSectionInputs["name"]].Value(),
					Duration:       d,
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
	cmds = append(cmds, m.UpdateFocus(msg))
	return m, tea.Batch(cmds...)
}

func (m CreateSection) View() string {
	var sb strings.Builder

	sb.WriteString(
		CenteredContainerStyle.Width(m.Width).
			Render(TitleStyle.Render("Create a Section")),
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
