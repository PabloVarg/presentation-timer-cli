package cli

import (
	"fmt"
	"os"

	"github.com/PabloVarg/presentation-timer-cli/internal/api"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var keyMap = []key.Binding{
	key.NewBinding(
		key.WithKeys("a"),
		key.WithHelp("a", "new presentation"),
	),
}

type ListPresentations struct {
	ProgramModel
	APIModel
	StyledComponent
	ListModel[api.Presentation]
}

func NewListPresentations(m ProgramModel) ListPresentations {
	l := NewDefaultList(NewDefaultDelegate(), keyMap)

	return ListPresentations{
		ProgramModel: m,
		ListModel: ListModel[api.Presentation]{
			list:     &l,
			itemizer: PresentationItemizer,
		},
		StyledComponent: StyledComponent{
			styles: map[string]lipgloss.Style{
				"list": containerStyle,
			},
		},
	}
}

func (m ListPresentations) Init() tea.Cmd {
	m.list.SetSize(
		m.width-m.styles["list"].GetHorizontalFrameSize(),
		m.height-m.styles["list"].GetVerticalFrameSize(),
	)

	return tea.Batch(func() tea.Msg {
		result, err := api.GetPresentations(m.api, os.LookupEnv)
		if err != nil {
			return fmt.Errorf("error loading data %s", err)
		}

		return result
	})
}

func (m ListPresentations) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := make([]tea.Cmd, 0)
	cmds = append(cmds, m.ProgramModel.Update(msg))

	switch msg := msg.(type) {
	case error:
		cmds = append(cmds, m.handleError(msg))
	case []api.Presentation:
		cmds = append(cmds, m.handleItems(msg...))
	case tea.WindowSizeMsg:
		m.list.SetSize(
			m.width-m.styles["list"].GetHorizontalFrameSize(),
			m.height-m.styles["list"].GetVerticalFrameSize(),
		)
	case tea.KeyMsg:
		switch msg.String() {
		case "a":
			nextModel := NewCreatePresentation(m.ProgramModel)
			return nextModel, nextModel.Init()
		}
	}

	model, cmd := m.list.Update(msg)
	m.list = &model
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m ListPresentations) View() string {
	return m.styles["list"].Render(m.list.View())
}
