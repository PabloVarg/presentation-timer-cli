package cli

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/PabloVarg/presentation-timer-cli/internal/api"
	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gorilla/websocket"
)

type RunPresentation struct {
	ProgramModel
	APIModel
	StyledComponent
	progress progress.Model

	presentationID int
	ctx            context.Context
	cancelCtx      context.CancelFunc
	websocketConn  *websocket.Conn
	readChannel    <-chan api.RunStatusResponse
	sendChannel    chan<- api.RunInput
	lastStatus     api.RunStatusResponse
}

type connectWsMsg struct {
	conn *websocket.Conn
	rxCh <-chan api.RunStatusResponse
	txCh chan<- api.RunInput
	err  error
}

func NewRunPresentation(m ProgramModel, ID int) RunPresentation {
	ctx, cancel := context.WithCancel(context.Background())

	return RunPresentation{
		ProgramModel:   m,
		presentationID: ID,
		ctx:            ctx,
		cancelCtx:      cancel,
		progress: progress.New(
			progress.WithWidth(m.Width-3*2),
			progress.WithSolidFill(string(Blue)),
			progress.WithoutPercentage(),
		),
		StyledComponent: StyledComponent{
			Styles: map[string]lipgloss.Style{
				"clock": lipgloss.NewStyle().
					Width(m.Width).
					AlignHorizontal(lipgloss.Center).
					Padding(5, 0).
					Bold(true).
					Blink(true),
				"progress": lipgloss.NewStyle().
					Padding(2, 3).Width(m.Width).AlignHorizontal(lipgloss.Center),
			},
		},
	}
}

func (m RunPresentation) Init() tea.Cmd {
	cmds := make([]tea.Cmd, 0, 2)
	cmds = append(cmds, Tick(100*time.Millisecond))

	ch, send, conn, err := api.ConnectToRun(m.ctx, m.Api, m.Logger, m.presentationID)
	if err != nil {
		cmds = append(cmds, func() tea.Msg {
			return ConnClosed
		})
		return tea.Batch(cmds...)
	}

	cmds = append(cmds, func() tea.Msg {
		return connectWsMsg{
			conn: conn,
			rxCh: ch,
			txCh: send,
			err:  err,
		}
	})

	return tea.Batch(cmds...)
}

func Tick(d time.Duration) tea.Cmd {
	return tea.Every(d, func(t time.Time) tea.Msg { return d })
}

func (m RunPresentation) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	m.Logger.Info("Update", "msg", msg)
	cmds := make([]tea.Cmd, 0)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.progress.Width = msg.Width - m.Styles["progress"].GetHorizontalPadding()*2
	case connectWsMsg:
		m.websocketConn = msg.conn
		m.sendChannel = msg.txCh
		m.readChannel = msg.rxCh

		m.sendChannel <- api.RunInput{
			Action: "status",
		}

		cmds = append(cmds, m.WaitForMessage())
	case api.RunStatusResponse:
		m.Logger.Info("update handling status response")
		m.lastStatus = msg
		cmds = append(cmds, m.WaitForMessage())
		m.Logger.Info("finish: update handling status response")
	case time.Duration:
		m.Logger.Info("tick")
		cmds = append(cmds, Tick(msg))
		if m.lastStatus.State != "running" || m.lastStatus.MsLeft <= 0 {
			break
		}

		m.lastStatus.MsLeft -= min(msg.Milliseconds(), m.lastStatus.MsLeft)
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc":
			if err := m.websocketConn.Close(); err != nil {
				m.Logger.Info("error closing ws")
			}
			m.cancelCtx()
			return Transition(NewListPresentations(m.ProgramModel))
		case "s":
			m.Logger.Info("reached before send")
			m.sendChannel <- api.RunInput{
				Action: "start",
			}
			m.Logger.Info("reached after send")
		case "t":
			if m.lastStatus.State == "running" {
				m.sendChannel <- api.RunInput{
					Action: "pause",
				}
				break
			}

			m.sendChannel <- api.RunInput{
				Action: "resume",
			}
		case "p":
			prevStep := int32(m.lastStatus.Step.Position - 2)
			m.sendChannel <- api.RunInput{
				Action: "step",
				Step:   &prevStep,
			}
		case "n":
			nextStep := int32(m.lastStatus.Step.Position)
			m.sendChannel <- api.RunInput{
				Action: "step",
				Step:   &nextStep,
			}
		}
		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit
		}
	case error:
		switch {
		case errors.Is(msg, ConnClosed):
			m.Logger.Info("received", "conn", "closed")
		}
	}

	return m, tea.Batch(cmds...)
}

func (m RunPresentation) View() string {
	var sb strings.Builder

	sb.WriteString(
		TitleStyle.Width(m.Width).
			AlignHorizontal(lipgloss.Center).
			Padding(2).
			Render(m.lastStatus.Step.Name),
	)
	sb.WriteString("\n")
	sb.WriteString(
		m.Styles["clock"].Render(
			fmt.Sprintf(
				"%02d:%02d",
				m.lastStatus.MsLeft/1_000/1_000,
				m.lastStatus.MsLeft%(1_000*1_000)/1_000,
			)),
	)
	sb.WriteString("\n")
	sb.WriteString(
		m.Styles["progress"].Render(
			m.progress.ViewAs(
				1 - float64(m.lastStatus.MsLeft)/float64(m.lastStatus.Step.Duration.Milliseconds()),
			),
		),
	)
	sb.WriteString("\n")

	return sb.String()
}

func (m RunPresentation) WaitForMessage() tea.Cmd {
	return func() tea.Msg {
		select {
		case <-m.ctx.Done():
			return nil
		case res, ok := <-m.readChannel:
			if !ok {
				return ConnClosed
			}

			return res
		}
	}
}
