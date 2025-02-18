package main

import (
	"fmt"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var counterStyle = lipgloss.NewStyle()

type model struct {
	buttons    buttons
	order      []position
	orderIndex int
	count      int
	playback   bool
	busy       bool
}

type waitMsg struct{}

func wait() tea.Msg {
	time.Sleep(time.Second)
	return waitMsg{}
}

func shortWait() tea.Msg {
	time.Sleep(time.Second / 3)
	return waitMsg{}
}

func (m model) Init() tea.Cmd {
	return tea.Sequence(
		wait,
		nextPosition,
	)
}

func (m model) handleInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	key := msg.String()
	if key == "ctrl+c" || key == "ctrl+d" {
		return m, tea.Quit
	}

	if m.busy {
		return m, nil
	}

	var pressedPosition position

	switch msg.String() {
	case "up":
		pressedPosition = top
	case "down":
		pressedPosition = bottom
	case "left":
		pressedPosition = left
	case "right":
		pressedPosition = right
	default:
		return m, nil
	}

	if m.order[m.orderIndex] != pressedPosition {
		m.orderIndex = 0
		m.order = []position{}
		m.count = 0
		m.buttons.Error()
		return m, tea.Sequence(
			blink(nil),
			wait,
			nextPosition,
		)

	}

	m.buttons.Button(pressedPosition).Blink()
	m.orderIndex++

	if m.orderIndex == len(m.order) {
		m.orderIndex = 0
		m.count++
		return m, tea.Sequence(
			blink(m.buttons.Button(pressedPosition)),
			wait,
			nextPosition,
		)
	}

	return m, blink(m.buttons.Button(pressedPosition))
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		buttonStyle = buttonStyle.Height(msg.Height / 3).Width(msg.Width / 3)
		counterStyle = buttonStyle.Align(lipgloss.Center).AlignVertical(lipgloss.Center)
		m.buttons.Button(top).style = buttonStyle.
			MarginLeft(((msg.Width / 2) - buttonStyle.GetWidth()/2) - 1)
		m.buttons.Button(top).ResetColor()

		m.buttons.Button(right).style = buttonStyle
		m.buttons.Button(right).ResetColor()

		m.buttons.Button(bottom).style = buttonStyle.
			MarginLeft(((msg.Width / 2) - buttonStyle.GetWidth()/2) - 1)
		m.buttons.Button(bottom).ResetColor()

		m.buttons.Button(left).style = buttonStyle
		m.buttons.Button(left).ResetColor()

	case tea.KeyMsg:
		return m.handleInput(msg)
	case blinkMsg:
		if msg.Button == nil {
			m.buttons.ResetColor()
		} else {
			msg.Button.ResetColor()
		}
		m.busy = false
	case position:
		m.busy = true
		m.order = append(m.order, msg)
		m.orderIndex = 0

		return m, playback
	case playbackMsg:
		// Last item of playback
		if m.orderIndex == len(m.order) {
			m.busy = false
			m.orderIndex = 0
			return m, nil
		}

		btn := m.buttons.Button(m.order[m.orderIndex])

		btn.Blink()
		m.orderIndex++

		return m, tea.Sequence(blink(btn), shortWait, playback)
	}

	return m, nil
}

func (m model) View() string {
	return lipgloss.JoinVertical(lipgloss.Top,
		m.buttons.Button(top).Render(),
		lipgloss.JoinHorizontal(
			lipgloss.Left,
			m.buttons.Button(left).Render(),
			counterStyle.Render(fmt.Sprintf("%d", m.count)),
			m.buttons.Button(right).Render(),
		),
		m.buttons.Button(bottom).Render(),
	)
}

func main() {
	buttons := newButtons()
	p := tea.NewProgram(
		model{buttons: buttons},
		tea.WithAltScreen(),
	)
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
