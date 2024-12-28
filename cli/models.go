package cli

import (
	"time"

	"github.com/PabloVarg/presentation-timer-cli/internal/api"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

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

type FormModel struct {
	focusIndex int
	inputs     []textinput.Model
}

func (f *FormModel) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(f.inputs))

	for i := range f.inputs {
		f.inputs[i], cmds[i] = f.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}
