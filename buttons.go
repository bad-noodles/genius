package main

import (
	"math/rand"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	sucessColor = lipgloss.Color("#fff")
	errorColor  = lipgloss.Color("#d11141")
	rightColor  = lipgloss.Color("#f37735")
	topColor    = lipgloss.Color("#00b159")
	bottomColor = lipgloss.Color("#00aedb")
	leftColor   = lipgloss.Color("#ffc425")
	buttonStyle = lipgloss.NewStyle().Width(10).Height(10)
)

type position int

const (
	top position = iota
	bottom
	left
	right
)

type button struct {
	style lipgloss.Style
	color lipgloss.Color
}

func (b button) Render() string {
	return b.style.Render()
}

func (b *button) Blink() {
	b.style = b.style.Background(sucessColor)
}

func (b *button) Error() {
	b.style = b.style.Background(errorColor)
}

func (b *button) ResetColor() {
	b.style = b.style.Background(b.color)
}

type buttons []*button

func (b buttons) Error() {
	for _, butt := range b {
		butt.Error()
	}
}

func (b buttons) ResetColor() {
	for _, butt := range b {
		butt.ResetColor()
	}
}

func (b buttons) Button(pos position) *button {
	return b[pos]
}

func newButtons() buttons {
	return buttons{
		&button{buttonStyle, topColor},
		&button{buttonStyle, bottomColor},
		&button{buttonStyle, leftColor},
		&button{buttonStyle, rightColor},
	}
}

func nextPosition() tea.Msg {
	return []position{top, right, bottom, left}[rand.Intn(4)]
}

type blinkMsg struct {
	Button *button
}

func blink(btn *button) tea.Cmd {
	return func() tea.Msg {
		time.Sleep(time.Second / 2)
		return blinkMsg{btn}
	}
}

type playbackMsg struct{}

func playback() tea.Msg {
	return playbackMsg{}
}

func blinkButton(b *button) tea.Cmd {
	return func() tea.Msg {
		b.Blink()
		return waitMsg{}
	}
}
