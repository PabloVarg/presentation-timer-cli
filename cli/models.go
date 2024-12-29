package cli

import (
	"log/slog"
	"time"

	"github.com/PabloVarg/presentation-timer-cli/internal/api"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ProgramModel struct {
	Logger *slog.Logger
	height int
	width  int
}

func (m *ProgramModel) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.width = msg.Width
	}

	return nil
}

type StyledComponent struct {
	styles map[string]lipgloss.Style
}

type APIModel struct {
	api api.APIClient
}

type ListModel[T any] struct {
	list     *list.Model
	itemizer func([]T) []list.Item
}

func (l *ListModel[T]) handleError(err error) tea.Cmd {
	previousLifetime := l.list.StatusMessageLifetime

	l.list.StatusMessageLifetime = 10 * time.Second
	cmd := l.list.NewStatusMessage(err.Error())
	l.list.StatusMessageLifetime = previousLifetime

	return cmd
}

func (l *ListModel[T]) handleItems(items ...T) tea.Cmd {
	cmds := make([]tea.Cmd, 0)

	cmds = append(cmds, l.list.SetItems(l.itemizer(items)))

	return tea.Batch(cmds...)
}

type (
	FormError = error
	FormModel struct {
		focusIndex int
		inputs     []textinput.Model
		err        FormError
	}
)

func (f *FormModel) UpdateForm(msg tea.Msg, sendKey tea.KeyType) {
	switch msg := msg.(type) {
	case FormError:
		f.err = msg
	case tea.KeyMsg:
		switch msg.Type {
		default:
			f.err = nil
		}
	}
}

func (f *FormModel) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(f.inputs))

	for i := range f.inputs {
		f.inputs[i], cmds[i] = f.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func transition(to tea.Model) (tea.Model, tea.Cmd) {
	return to, to.Init()
}
