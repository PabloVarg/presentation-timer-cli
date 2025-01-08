package cli

import (
	"os"
	"strings"

	"github.com/PabloVarg/presentation-timer-cli/internal/api"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type EditPresentation struct {
	ProgramModel
	APIModel
	StyledComponent
	FormModel

	isLoading      bool
	PresentationID int
}

var editPresentationInputs = map[string]int{
	"name": 0,
}

func NewEditPresentation(m ProgramModel, ID int) EditPresentation {
	nameInput := NewDefaultTextInput()
	nameInput.Placeholder = "My Presentation"
	nameInput.Prompt = "Name: "

	nameInput.Focus()

	return EditPresentation{
		PresentationID: ID,
		ProgramModel:   m,
		FormModel: FormModel{
			inputs: []textinput.Model{
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
			p, err := api.GetPresentation(m.api, os.LookupEnv, m.PresentationID)
			if err != nil {
				return FetchError{
					err: err.Error(),
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
		m.inputs[editPresentationInputs["name"]].SetValue(msg.Name)
		m.isLoading = false
	case FetchError:
		m.inputs = nil
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit
		case tea.KeyEsc:
			return transition(NewListPresentations(m.ProgramModel))
		case tea.KeyEnter:
			cmds = append(cmds, func() tea.Msg {
				err := api.UpdatePresentation(m.api, os.LookupEnv, api.EditPresentationMsg{
					ID:   m.PresentationID,
					Name: m.inputs[editPresentationInputs["name"]].Value(),
				})
				if err != nil {
					return FormError{
						err: err.Error(),
					}
				}
				return tea.KeyMsg{
					Type: tea.KeyEsc,
				}
			})
		}
	}

	m.UpdateForm(msg, tea.KeyEnter)

	cmds = append(cmds, m.updateInputs(msg))
	return m, tea.Batch(cmds...)
}

func (m EditPresentation) View() string {
	var sb strings.Builder

	sb.WriteString(
		centeredContainerStyle.Width(m.width).
			Render(titleStyle.Render("Edit Presentation")),
	)
	sb.WriteRune('\n')

	for i := range m.inputs {
		sb.WriteString(m.inputs[i].View())
		sb.WriteRune('\n')
	}

	if m.err != nil {
		sb.WriteRune('\n')
		sb.WriteString(errorStyle.Render(m.err.Error()))
	}

	return containerStyle.Render(sb.String())
}
