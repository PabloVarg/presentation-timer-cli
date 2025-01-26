package cli

import (
	"os"
	"strconv"
	"strings"

	"github.com/PabloVarg/presentation-timer-cli/internal/api"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type MoveSection struct {
	ProgramModel
	APIModel
	StyledComponent
	FormModel

	presentationID int
	sectionID      int
}

var moveSectionInputs = map[string]int{
	"move": 0,
}

func NewMoveSection(m ProgramModel, presentationID int, sectionID int) MoveSection {
	nameInput := NewDefaultTextInput()
	nameInput.Placeholder = "+3"
	nameInput.Prompt = "Move: "

	nameInput.Focus()

	return MoveSection{
		sectionID:      sectionID,
		presentationID: presentationID,
		ProgramModel:   m,
		FormModel: FormModel{
			Inputs: []textinput.Model{
				nameInput,
			},
		},
	}
}

func (m MoveSection) Init() tea.Cmd {
	return textinput.Blink
}

func (m MoveSection) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
				move, err := strconv.Atoi(m.Inputs[moveSectionInputs["move"]].Value())
				if err != nil {
					return FormError{
						Err: err.Error(),
					}
				}

				err = api.MoveSection(m.Api, os.LookupEnv, api.MoveSectionMsg{
					SectionID: m.sectionID,
					Move:      move,
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

func (m MoveSection) View() string {
	var sb strings.Builder

	sb.WriteString(
		CenteredContainerStyle.Width(m.Width).
			Render(TitleStyle.Render("Move Section")),
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
