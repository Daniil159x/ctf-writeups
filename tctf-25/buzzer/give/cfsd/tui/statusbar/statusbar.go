package statusbar

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type tickMsg struct{}

type NewStatus struct {
	Msg string
}

type StatusBar struct {
	style lipgloss.Style

	Deadline time.Time
	Message  string
}

func NewStatusBar(deadline time.Time) *StatusBar {
	return &StatusBar{
		style:    lipgloss.NewStyle().Border(lipgloss.NormalBorder()),
		Deadline: deadline,
		Message:  "no new message",
	}
}

func doTick() tea.Cmd {
	return tea.Tick(time.Millisecond*500, func(t time.Time) tea.Msg {
		return tickMsg{}
	})
}

func (s *StatusBar) SetSize(w int) {
	s.style = s.style.Width(w - s.style.GetHorizontalBorderSize())
}

// Init implements tea.Model.
func (s *StatusBar) Init() tea.Cmd {
	return doTick()
}

// Update implements tea.Model.
func (s *StatusBar) Update(msg tea.Msg) (*StatusBar, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		s.SetSize(msg.Width)
	case tickMsg:
		cmd = doTick()
	case NewStatus:
		s.Message = msg.Msg
	}

	return s, cmd
}

func (s *StatusBar) getRemained() string {
	r := time.Until(s.Deadline)
	minutes := int(r.Minutes())
	seconds := int(r.Seconds()) % 60
	milliseconds := int(r.Milliseconds()) % 1000

	return fmt.Sprintf("Осталось: %02d:%02d.%03d", minutes, seconds, milliseconds)
}

// View implements tea.Model.
func (s *StatusBar) View() string {
	remained := s.getRemained()

	left := lipgloss.Width(s.Message)
	right := lipgloss.Width(remained)
	w := s.style.GetWidth()

	gap := w - left - right - 1
	remained = lipgloss.NewStyle().MarginLeft(gap).Render(remained)

	return s.style.Render(s.Message, remained)
}
