package cli

import (
	"fmt"
	"os"

	"github.com/PabloVarg/presentation-timer-cli/internal/api"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ListPresentations struct {
	APIModel
	StyledComponent
	ListModel[api.Presentation]
}

func NewListPresentations() ListPresentations {
	l := NewDefaultList(NewDefaultDelegate())

	return ListPresentations{
		ListModel: ListModel[api.Presentation]{
			list:     &l,
			itemizer: PresentationItemizer,
		},
		StyledComponent: StyledComponent{
			styles: map[string]lipgloss.Style{
				"list": lipgloss.
					NewStyle().
					Padding(1),
			},
		},
	}
}

func (m ListPresentations) Init() tea.Cmd {
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

	switch msg := msg.(type) {
	case error:
		cmds = append(cmds, m.handleError(msg))
	case []api.Presentation:
		cmds = append(cmds, m.handleItems(msg...))
	case tea.WindowSizeMsg:
		m.list.SetSize(
			msg.Width-m.styles["list"].GetHorizontalFrameSize(),
			msg.Height-m.styles["list"].GetVerticalFrameSize(),
		)
	}

	model, cmd := m.list.Update(msg)
	m.list = &model
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m ListPresentations) View() string {
	return m.styles["list"].Render(m.list.View())
}
