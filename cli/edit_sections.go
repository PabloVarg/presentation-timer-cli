package cli

import (
	"os"
	"strings"
	"time"

	"github.com/PabloVarg/presentation-timer-cli/internal/api"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type EditSection struct {
	ProgramModel
	APIModel
	StyledComponent
	FormModel

	sectionID      int
	presentationID int
}

var editSectionInputs = map[string]int{
	"name":     0,
	"duration": 1,
}

func NewEditSection(m ProgramModel, presentationID int, sectionID int) EditSection {
	nameInput := NewDefaultTextInput()
	nameInput.Placeholder = "My Section"
	nameInput.Prompt = "Name: "

	durationInput := NewDefaultTextInput()
	durationInput.Placeholder = "15m20s"
	durationInput.Prompt = "Duration: "

	nameInput.Focus()

	return EditSection{
		presentationID: presentationID,
		sectionID:      sectionID,
		ProgramModel:   m,
		FormModel: FormModel{
			Inputs: []textinput.Model{
				nameInput,
				durationInput,
			},
		},
	}
}

func (m EditSection) Init() tea.Cmd {
	return tea.Batch(
		textinput.Blink,
		func() tea.Msg {
			p, err := api.GetSection(m.Api, os.LookupEnv, m.sectionID)
			if err != nil {
				return FetchError{
					Err: err.Error(),
				}
			}

			return p
		},
	)
}

func (m EditSection) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := make([]tea.Cmd, 0)
	cmds = append(cmds, m.ProgramModel.Update(msg))

	switch msg := msg.(type) {
	case api.Section:
		m.Inputs[editSectionInputs["name"]].SetValue(msg.Name)
		m.Inputs[editSectionInputs["duration"]].SetValue(msg.Duration.String())
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

				err = api.UpdateSection(m.Api, os.LookupEnv, api.EditSectionMsg{
					ID:       m.sectionID,
					Name:     m.Inputs[createSectionInputs["name"]].Value(),
					Duration: d,
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

func (m EditSection) View() string {
	var sb strings.Builder

	sb.WriteString(
		CenteredContainerStyle.Width(m.Width).
			Render(TitleStyle.Render("Edit Section")),
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
