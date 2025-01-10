package cli

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ConfirmationModel struct {
	ProgramModel
	StyledComponent
	parentModel tea.Model
	parentCmds  []tea.Cmd
	acceptFunc  func() tea.Cmd
	rejectFunc  func() tea.Cmd
}

func NewConfirmationModel(
	parentModel tea.Model,
	accept func() tea.Cmd,
	opts ...func(m *ConfirmationModel),
) ConfirmationModel {
	m := ConfirmationModel{
		parentModel: parentModel,
		acceptFunc:  accept,
	}

	for _, opt := range opts {
		opt(&m)
	}

	m.StyledComponent = StyledComponent{
		Styles: map[string]lipgloss.Style{
			"layout": lipgloss.NewStyle().
				Width(m.ProgramModel.Width).
				Height(m.ProgramModel.Height).
				Align(lipgloss.Center, lipgloss.Bottom),
			"container": lipgloss.NewStyle().
				Width(m.ProgramModel.Width/2).
				Border(lipgloss.RoundedBorder(), true).
				BorderForeground(Red).
				Align(lipgloss.Center).
				MarginBottom(5).
				Padding(1),
			"content": lipgloss.NewStyle().PaddingTop(1),
		},
	}

	return m
}

func WithRejectFunc(f func() tea.Cmd) func(m *ConfirmationModel) {
	return func(m *ConfirmationModel) {
		m.rejectFunc = f
	}
}

func WithProgramModel(pm ProgramModel) func(m *ConfirmationModel) {
	return func(m *ConfirmationModel) {
		m.ProgramModel = pm
	}
}

func (m ConfirmationModel) Init() tea.Cmd {
	return nil
}

func (m ConfirmationModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "y", "Y", "enter":
			return m.parentModel, m.acceptFunc()
		case "n", "N", "esc", "q":
			f := func() tea.Cmd { return nil }
			if m.rejectFunc != nil {
				f = m.rejectFunc
			}

			return m.parentModel, f()
		}
		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m ConfirmationModel) View() string {
	var sb strings.Builder

	sb.WriteString(TitleStyle.Render("Confirm action\n"))
	sb.WriteString(m.Styles["content"].Render("(Y)es (N)o"))

	return m.Styles["layout"].Render(m.Styles["container"].Render(sb.String()))
}
