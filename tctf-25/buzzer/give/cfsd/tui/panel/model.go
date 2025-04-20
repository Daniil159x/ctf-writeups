package panel

import (
	"cfsd/internal/command"
	"cfsd/pkg/sugar"
	"cfsd/tui/cmdselector"
	"cfsd/tui/contentview"
	"cfsd/tui/fixhelp"
	"cfsd/tui/statusbar"
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	formWidget = iota
	viewWidget
	maxWidget

	MaxWidth  = 200
	MaxHeight = 100
)

type Model struct {
	top    *cmdselector.CmdSelector
	middle *contentview.ContentView
	bottom *statusbar.StatusBar

	currentWidget int

	width  int
	height int

	ThisProgram *tea.Program
}

func NewModel() (*Model, error) {
	deadline := time.Now().Add(time.Minute * 10)

	middle, err := contentview.NewContentView()
	if err != nil {
		return nil, fmt.Errorf("failed to create viewport: %w", err)
	}

	return &Model{
		top:           cmdselector.NewSelect(),
		middle:        middle,
		bottom:        statusbar.NewStatusBar(deadline),
		currentWidget: formWidget,
	}, nil
}

func (m *Model) Init() tea.Cmd {
	m.top.ThisProgram = m.ThisProgram
	return tea.Batch(
		m.top.Init(),
		m.middle.Init(),
		m.bottom.Init(),
	)
}

func (m *Model) handleResize(msg tea.WindowSizeMsg) {
	m.width = min(MaxWidth, msg.Width) - 4
	m.height = min(MaxHeight, msg.Height) - 4

	m.top.SetWidth(m.width)
	m.bottom.SetSize(m.width)

	h := m.height -
		lipgloss.Height(m.bottom.View()) -
		lipgloss.Height(m.top.View())

	h = max(10, h)
	m.middle.SetSize(m.width, h)
}

func (m *Model) changeFocus(widget int) {
	m.currentWidget = widget
	m.top.Focus = (m.currentWidget == formWidget)
	m.middle.Focus = (m.currentWidget == viewWidget)
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	batch := make([]tea.Cmd, 0)

	switch msgTyped := msg.(type) {
	case tea.KeyMsg:

		switch {
		case key.Matches(msgTyped, fixhelp.ExitKey):
			return m, tea.Quit
		case key.Matches(msgTyped, fixhelp.NextWidgetKey):
			m.changeFocus((m.currentWidget + 1) % maxWidget)
			return m, nil
		}

		if m.currentWidget == formWidget {
			sugar.Update2(&m.top, msg, &batch)
		} else {
			sugar.Update2(&m.middle, msg, &batch)
		}
		return m, tea.Batch(batch...)

	case tea.WindowSizeMsg:
		m.handleResize(msgTyped)
		return m, tea.Batch(batch...)

	case cmdselector.CmdSelected:
		// form already show list file
		if msgTyped.Command != cmdselector.CmdListFile {
			m.changeFocus(viewWidget)
			return m, m.handleCmd(msgTyped)
		}

	case CmdFinished:
		m.middle.SetContent(msgTyped.output, msgTyped.cmd.Arg)

		m.bottom.Message = fmt.Sprintf("команда завершена: err=%v", msgTyped.err)
		return m, nil

	case command.CmdError:
		m.bottom.Message = fmt.Sprintf("команда зафейлилась: err=%v", msgTyped.Err)
		m.middle.SetContent(msgTyped.Stderr, "stderr")

		return m, nil
	}

	sugar.Update2(&m.middle, msg, &batch)
	sugar.Update2(&m.top, msg, &batch)
	sugar.Update2(&m.bottom, msg, &batch)

	return m, tea.Batch(batch...)
}

func (m *Model) View() string {
	return lipgloss.NewStyle().
		Border(lipgloss.DoubleBorder()).
		Render(
			lipgloss.JoinVertical(
				lipgloss.Left,
				m.top.View(),
				m.middle.View(),
				m.bottom.View(),
			),
		)
}
