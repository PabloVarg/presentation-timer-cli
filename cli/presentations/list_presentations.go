package presentations

import (
	"fmt"
	"os"

	"github.com/PabloVarg/presentation-timer-cli/cli"
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
	key.NewBinding(
		key.WithKeys("D"),
		key.WithHelp("D", "delete presentation"),
	),
	key.NewBinding(
		key.WithKeys("R"),
		key.WithHelp("R", "refresh"),
	),
}

type ListPresentations struct {
	cli.ProgramModel
	cli.APIModel
	cli.StyledComponent
	cli.ListModel[api.Presentation]
}

func NewListPresentations(m cli.ProgramModel) ListPresentations {
	l := cli.NewDefaultList(cli.NewDefaultDelegate(), keyMap, "Presentations")

	return ListPresentations{
		ProgramModel: m,
		ListModel: cli.ListModel[api.Presentation]{
			List:     &l,
			Itemizer: cli.PresentationItemizer,
		},
		StyledComponent: cli.StyledComponent{
			Styles: map[string]lipgloss.Style{
				"list": cli.ContainerStyle,
			},
		},
	}
}

func (m ListPresentations) Init() tea.Cmd {
	m.List.SetSize(
		m.Width-m.Styles["list"].GetHorizontalFrameSize(),
		m.Height-m.Styles["list"].GetVerticalFrameSize(),
	)

	return tea.Batch(m.retrievePresentations())
}

func (m ListPresentations) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := make([]tea.Cmd, 0)
	cmds = append(cmds, m.ProgramModel.Update(msg))

	switch msg := msg.(type) {
	case error:
		cmds = append(cmds, m.HandleError(msg))
	case []api.Presentation:
		cmds = append(cmds, m.HandleItems(msg...))
	case tea.WindowSizeMsg:
		m.List.SetSize(
			m.Width-m.Styles["list"].GetHorizontalFrameSize(),
			m.Height-m.Styles["list"].GetVerticalFrameSize(),
		)
	case tea.KeyMsg:
		switch msg.String() {
		case "a":
			return cli.Transition(NewCreatePresentation(m.ProgramModel))
		case "D":
			nextModel := cli.NewConfirmationModel(m, m.deleteSelectedItem, cli.WithProgramModel(m.ProgramModel))
			return nextModel, nextModel.Init()
		case "c":
			item, ok := m.List.SelectedItem().(cli.PresentationItem)
			if !ok {
				panic("received unexpected value type")
			}

			return cli.Transition(NewEditPresentation(m.ProgramModel, item.ID))
		case "R":
			cmds = append(cmds, m.retrievePresentations())
		}
	}

	model, cmd := m.List.Update(msg)
	m.List = &model
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m ListPresentations) View() string {
	return m.Styles["list"].Render(m.List.View())
}

func (m ListPresentations) retrievePresentations() tea.Cmd {
	return func() tea.Msg {
		result, err := api.GetPresentations(m.Api, os.LookupEnv)
		if err != nil {
			return fmt.Errorf("error loading data %s", err)
		}

		return result
	}
}

func (m ListPresentations) deleteSelectedItem() tea.Cmd {
	return tea.Sequence(
		func() tea.Msg {
			selectedItem, ok := m.List.SelectedItem().(cli.PresentationItem)
			if !ok {
				panic("received unexpected value type")
			}
			err := api.DeletePresentation(m.Api, os.LookupEnv, selectedItem.ID)
			if err != nil {
				return fmt.Errorf("could not delete %s", selectedItem.Name)
			}

			return nil
		},
		m.retrievePresentations(),
	)
}
