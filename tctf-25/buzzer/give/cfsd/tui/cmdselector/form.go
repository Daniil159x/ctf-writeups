package cmdselector

import (
	"bytes"
	"cfsd/internal/command"
	"cfsd/pkg/sugar"
	"cfsd/tui/fixhelp"

	"github.com/charmbracelet/bubbles/help"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

const (
	CmdListFile    = "Список файлов"
	CmdOpenFile    = "Открыть файл"
	CmdOpenShell   = "Открыть терминал"
	CmdChangeTheme = "Поменять тему панели"

	MaxHeight = 8
)

type CmdSelected struct {
	Command string
	Arg     string
}

type CmdSelector struct {
	value CmdSelected
	form  *huh.Form
	help  help.Model

	width  int
	height int

	Focus       bool
	ThisProgram *tea.Program
}

func (m *CmdSelector) SetWidth(width int) {
	m.form = m.form.WithWidth(width)
	m.help.Width = width
	m.width = width
}

func (m *CmdSelector) SetHeight(height int) {
	m.height = min(MaxHeight, height)
	m.form = m.form.WithHeight(height)
}

func NewSelect() *CmdSelector {
	m := &CmdSelector{
		Focus: true,
		help:  help.New(),
	}
	m.form = m.newForm()
	return m
}

func (m *CmdSelector) newForm() *huh.Form {
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Выбери действие").
				Options(huh.NewOptions(
					CmdListFile,
					CmdOpenFile,
					CmdOpenShell,
					CmdChangeTheme,
				)...).
				Value(&m.value.Command),
		),
		huh.NewGroup(
			huh.NewSelect[string]().
				TitleFunc(func() string {
					switch m.value.Command {
					case CmdListFile:
						return "Вот все твои файлы"
					case CmdOpenFile:
						return "Выбери файл для открытия"
					case CmdChangeTheme:
						return "Выбери новую тему"
					case CmdOpenShell:
						return "Выбери тип оболочки"
					default:
						return "Ты как тут оказался?"
					}
				}, &m.value.Command).
				OptionsFunc(func() []huh.Option[string] {
					switch m.value.Command {
					case CmdListFile, CmdOpenFile:
						return huh.NewOptions(m.getFiles()...)
					case CmdChangeTheme:
						return huh.NewOptions("not impl")
					case CmdOpenShell:
						return huh.NewOptions("posix shell")
					}
					return nil
				}, &m.value.Command).
				Value(&m.value.Arg),
		),
	).
		WithHeight(MaxHeight).
		WithWidth(m.width).
		WithShowHelp(false) // draw ourselves

	form.SubmitCmd = m.submitCallback
	return form
}

func (m *CmdSelector) submitCallback() tea.Msg {
	return CmdSelected(m.value)
}

func (m *CmdSelector) getFiles() []string {
	cmd := command.MakeListFiles()

	var errContent bytes.Buffer
	cmd.Stderr = &errContent

	content, err := cmd.Output()

	if err != nil {
		m.ThisProgram.Send(command.CmdError{
			Stderr: errContent.Bytes(),
			Err:    err,
		})
	}

	bLines := bytes.Split(content, []byte("\n"))
	result := make([]string, 0, len(bLines))
	for _, bl := range bLines {
		line := string(bytes.TrimSpace(bl))
		if len(line) > 0 {
			result = append(result, line)
		}
	}

	return result
}

// Init implements tea.Model.
func (m *CmdSelector) Init() tea.Cmd {
	return m.form.Init()
}

// Update implements tea.Model.
func (m *CmdSelector) Update(msg tea.Msg) (*CmdSelector, tea.Cmd) {

	batch := make([]tea.Cmd, 0)

	sugar.Update(&m.form, msg, &batch)

	if m.form.State == huh.StateCompleted {
		m.form = m.newForm()
		batch = append(batch, m.form.Init())
	}

	return m, tea.Batch(batch...)
}

// View implements tea.Model.
func (m *CmdSelector) View() string {
	help := ""
	var border lipgloss.Border
	if m.Focus {
		help = m.help.View(fixhelp.FromBindings(m.form.KeyBinds()))
		border = lipgloss.ASCIIBorder()
	} else {
		border = lipgloss.HiddenBorder()
	}

	return lipgloss.NewStyle().
		Border(border).
		AlignVertical(lipgloss.Left).
		Render(
			m.form.View(),
			help,
		)
}
