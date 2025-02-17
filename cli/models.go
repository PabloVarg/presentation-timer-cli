package cli

import (
	"errors"
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
	Height int
	Width  int
}

func (m *ProgramModel) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Height = msg.Height
		m.Width = msg.Width
	}

	return nil
}

type StyledComponent struct {
	Styles map[string]lipgloss.Style
}

type APIModel struct {
	Api api.APIClient
}

type ListModel[T any] struct {
	List     *list.Model
	Itemizer func([]T) []list.Item
}

func (l *ListModel[T]) HandleError(err error) tea.Cmd {
	previousLifetime := l.List.StatusMessageLifetime

	l.List.StatusMessageLifetime = 10 * time.Second
	cmd := l.List.NewStatusMessage(err.Error())
	l.List.StatusMessageLifetime = previousLifetime

	return cmd
}

func (l *ListModel[T]) HandleItems(items ...T) tea.Cmd {
	cmds := make([]tea.Cmd, 0)

	cmds = append(cmds, l.List.SetItems(l.Itemizer(items)))

	return tea.Batch(cmds...)
}

type FormError struct {
	Err string
}

func (e FormError) Error() string {
	return e.Err
}

type FetchError struct {
	Err string
}

func (e FetchError) Error() string {
	return e.Err
}

type (
	FormModel struct {
		FocusIndex int
		Inputs     []textinput.Model
		Err        *FormError
	}
)

func (f *FormModel) UpdateForm(msg tea.Msg, sendKey tea.KeyType) {
	switch msg := msg.(type) {
	case FormError:
		f.Err = &msg
	case tea.KeyMsg:
		switch msg.Type {
		default:
			f.Err = nil
		}
	}
}

func (f *FormModel) UpdateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(f.Inputs))

	for i := range f.Inputs {
		f.Inputs[i], cmds[i] = f.Inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func (f *FormModel) UpdateFocus(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, 0, 1)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "tab", "shift+tab", "up", "down":
			s := msg.String()

			// Cycle indexes
			if s == "up" || s == "shift+tab" {
				f.FocusIndex--
			} else {
				f.FocusIndex++
			}

			if f.FocusIndex > len(f.Inputs) {
				f.FocusIndex = 0
			} else if f.FocusIndex < 0 {
				f.FocusIndex = len(f.Inputs)
			}

			for i := 0; i <= len(f.Inputs)-1; i++ {
				if i == f.FocusIndex {
					// Set focused state
					cmds = append(cmds, f.Inputs[i].Focus())
					continue
				}
				// Remove focused state
				f.Inputs[i].Blur()
			}
		}
	}

	return tea.Batch(cmds...)
}

func Transition(to tea.Model) (tea.Model, tea.Cmd) {
	return to, to.Init()
}

var ConnClosed = errors.New("connection lost")
