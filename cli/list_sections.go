package cli

import (
	"fmt"
	"os"

	"github.com/PabloVarg/presentation-timer-cli/internal/api"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var SectionsKeyMap = []key.Binding{
	key.NewBinding(
		key.WithKeys("a"),
		key.WithHelp("a", "new section"),
	),
	key.NewBinding(
		key.WithKeys("D"),
		key.WithHelp("D", "delete section"),
	),
	key.NewBinding(
		key.WithKeys("R"),
		key.WithHelp("R", "refresh"),
	),
}

type ListSections struct {
	ProgramModel
	APIModel
	StyledComponent
	ListModel[api.Section]

	presentationID int
}

func NewListSections(m ProgramModel, presentationID int) ListSections {
	l := NewDefaultList(NewDefaultDelegate(), SectionsKeyMap, "Sections")

	return ListSections{
		ProgramModel: m,
		ListModel: ListModel[api.Section]{
			List:     &l,
			Itemizer: SectionItemizer,
		},
		StyledComponent: StyledComponent{
			Styles: map[string]lipgloss.Style{
				"list": ContainerStyle,
			},
		},
		presentationID: presentationID,
	}
}

func (m ListSections) Init() tea.Cmd {
	m.List.SetSize(
		m.Width-m.Styles["list"].GetHorizontalFrameSize(),
		m.Height-m.Styles["list"].GetVerticalFrameSize(),
	)

	return tea.Batch(m.retrieveSections())
}

func (m ListSections) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := make([]tea.Cmd, 0)
	cmds = append(cmds, m.ProgramModel.Update(msg))

	switch msg := msg.(type) {
	case error:
		cmds = append(cmds, m.HandleError(msg))
	case []api.Section:
		cmds = append(cmds, m.HandleItems(msg...))
	case tea.WindowSizeMsg:
		m.List.SetSize(
			m.Width-m.Styles["list"].GetHorizontalFrameSize(),
			m.Height-m.Styles["list"].GetVerticalFrameSize(),
		)
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEsc:
			return Transition(NewListPresentations(m.ProgramModel))
		}

		if m.List.FilterState() == list.Filtering {
			break
		}
		switch msg.String() {
		case "a":
			return Transition(NewCreateSection(m.ProgramModel, m.presentationID))
		case "D":
			return Transition(NewConfirmationModel(m, m.deleteSelectedItem, WithProgramModel(m.ProgramModel)))
		case "c":
			item, ok := m.List.SelectedItem().(SectionItem)
			if !ok {
				panic("received unexpected value type")
			}

			return Transition(NewEditSection(m.ProgramModel, m.presentationID, item.ID))
		case "m":
			item, ok := m.List.SelectedItem().(SectionItem)
			if !ok {
				panic("received unexpected value type")
			}

			return Transition(NewMoveSection(m.ProgramModel, m.presentationID, item.ID))
		case "R":
			cmds = append(cmds, m.retrieveSections())
		}
	}

	model, cmd := m.List.Update(msg)
	m.List = &model
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m ListSections) View() string {
	return m.Styles["list"].Render(m.List.View())
}

func (m ListSections) retrieveSections() tea.Cmd {
	return func() tea.Msg {
		result, err := api.GetSections(m.Api, os.LookupEnv, m.presentationID)
		if err != nil {
			return fmt.Errorf("error loading data %s", err)
		}

		return result
	}
}

func (m ListSections) deleteSelectedItem() tea.Cmd {
	return tea.Sequence(
		func() tea.Msg {
			selectedItem, ok := m.List.SelectedItem().(SectionItem)
			if !ok {
				panic("received unexpected value type")
			}
			err := api.DeleteSection(m.Api, os.LookupEnv, selectedItem.ID)
			if err != nil {
				return fmt.Errorf("could not delete %s", selectedItem.Name)
			}

			return nil
		},
		m.retrieveSections(),
	)
}
