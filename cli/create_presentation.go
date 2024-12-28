package cli

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type CreatePresentation struct {
	StyledComponent
	FormModel
}

func NewCreatePresentation() CreatePresentation {
	nameInput := textinput.New()
	nameInput.Placeholder = "Name"

	nameInput.Focus()

	return CreatePresentation{
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

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			nextModel := NewListPresentations()
			return nextModel, nextModel.Init()
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
