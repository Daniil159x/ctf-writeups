package contentview

import (
	"bytes"
	"cfsd/pkg/sugar"
	"cfsd/tui/fixhelp"
	_ "embed"
	"fmt"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

//go:embed intro_stub.md
var stubMarkdown []byte

type ContentView struct {
	view     viewport.Model
	help     help.Model
	renderer render

	origin     []byte
	originMeta string

	needRerender bool
	counter      int

	Focus bool
}

func NewContentView() (*ContentView, error) {
	view := &ContentView{
		view:         viewport.New(40, 40),
		help:         help.New(),
		needRerender: false,
		Focus:        false,
		origin:       stubMarkdown,
		originMeta:   "hello.md",
	}

	return view, nil
}

func (c *ContentView) SetContent(content []byte, meta string) {
	if c.originMeta == meta && bytes.Equal(content, c.origin) {
		return
	}

	if c.originMeta != meta {
		c.renderer = nil
	}
	c.needRerender = true

	c.origin = content
	c.originMeta = meta

	if len(c.origin) == 0 {
		c.origin = stubMarkdown
		c.originMeta = "hello.md"
	}
}

func (c *ContentView) SetSize(width int, height int) {
	c.view.Height = height - 2 // status + help
	c.view.Width = width
	c.needRerender = true

	// NB: reset render for apply new geometry
	c.renderer = nil
}

// Init implements tea.Model.
func (c *ContentView) Init() tea.Cmd {
	return c.view.Init()
}

// Update implements tea.Model.
func (c *ContentView) Update(msg tea.Msg) (*ContentView, tea.Cmd) {
	batch := make([]tea.Cmd, 0)
	sugar.Update2(&c.view, msg, &batch)

	return c, tea.Batch(batch...)
}

func (c *ContentView) ScrollPercent() int {
	return int(c.view.ScrollPercent() * 100)
}

func (c *ContentView) rerender() {
	if !c.needRerender {
		return
	}
	c.needRerender = false

	// NB: apply new geometry
	var err error
	if c.renderer == nil {
		c.renderer, err = c.getMimeRender(Filename2Mime(c.originMeta))
	}

	var str string
	if err == nil {
		var resultBytes []byte

		resultBytes, err = c.renderer.RenderBytes(c.origin)
		if err == nil {
			str = string(resultBytes)
		}
	}

	if err != nil {
		str = fmt.Sprintf("failed to render content: %v", err)
	}

	c.view.SetContent(str)
	c.counter++
}

func (c *ContentView) View() string {
	c.rerender()

	help := ""
	var border lipgloss.Border
	if c.Focus {
		border = lipgloss.ASCIIBorder()
		help = c.help.View(fixhelp.FromViewport(c.view.KeyMap))
	} else {
		border = lipgloss.HiddenBorder()
	}

	scroll := fmt.Sprintf("%s\tscroll: %3d%%", c.originMeta, c.ScrollPercent())

	scroll = lipgloss.NewStyle().
		Width(c.view.Width - lipgloss.Width(help)).
		AlignHorizontal(lipgloss.Right).
		Italic(true).
		Render(scroll)

	footer := lipgloss.JoinHorizontal(lipgloss.Top, help, scroll)

	return lipgloss.NewStyle().
		Border(border).
		Render(
			lipgloss.JoinVertical(
				lipgloss.Left,
				c.view.View(),
				footer,
			),
		)
}
