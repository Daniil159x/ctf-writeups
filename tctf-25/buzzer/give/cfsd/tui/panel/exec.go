package panel

import (
	"bytes"
	"cfsd/internal/command"
	"cfsd/tui/cmdselector"
	"fmt"
	"os/exec"

	tea "github.com/charmbracelet/bubbletea"
)

type CmdStarted struct{}

type CmdNewLine struct {
	line string
}

type CmdFinished struct {
	cmd    cmdselector.CmdSelected
	output []byte
	err    error
}

func (m *Model) handleCmd(msg cmdselector.CmdSelected) tea.Cmd {
	switch msg.Command {
	case cmdselector.CmdOpenShell:
		return m.execInteractiveCommand(msg, command.MakeShell())
	case cmdselector.CmdOpenFile:
		return m.execNonInteractiveCommand(msg, command.MakeOpenFiles(msg.Arg))
	}
	return func() tea.Msg {
		return command.CmdError{
			Err: fmt.Errorf("Нет команды для '%s'", msg.Command),
		}
	}
}

func (m *Model) execInteractiveCommand(cmd cmdselector.CmdSelected, c *exec.Cmd) tea.Cmd {
	return tea.ExecProcess(c, func(err error) tea.Msg {
		return CmdFinished{
			err: err,
		}
	})
}

func (m *Model) execNonInteractiveCommand(cmd cmdselector.CmdSelected, c *exec.Cmd) tea.Cmd {
	return func() tea.Msg {

		var errContent bytes.Buffer
		c.Stderr = &errContent

		out, err := c.Output()
		if err != nil {
			return command.CmdError{
				Stderr: errContent.Bytes(),
				Err:    err,
			}
		}

		return CmdFinished{
			cmd:    cmd,
			output: out,
			err:    nil,
		}
	}
}
