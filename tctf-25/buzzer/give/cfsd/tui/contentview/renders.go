package contentview

import (
	"fmt"
	"path"

	"github.com/charmbracelet/glamour"
)

type ContentType int

const (
	SimpleText ContentType = iota
	Markdown
	Image
)

func Filename2Mime(filename string) ContentType {
	switch path.Ext(filename) {
	case ".md":
		return Markdown
	case ".png", ".jpeg", ".jpg":
		return Image
	default:
		return SimpleText
	}
}

type render interface {
	RenderBytes(in []byte) ([]byte, error)
}

func (c *ContentView) getMimeRender(mime ContentType) (render, error) {
	switch mime {
	case Markdown:
		return c.newMarkdownRender()
	case Image:
		return c.newImageRender()
	case SimpleText:
		return c.newSimpleTextRender()
	default:
		return nil, fmt.Errorf("not support mime type")
	}
}

func (c *ContentView) newMarkdownRender() (render, error) {
	renderer, err := glamour.NewTermRenderer(
		glamour.WithStandardStyle("dark"),
		glamour.WithWordWrap(c.view.Width),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create glamour render: %w", err)
	}
	return renderer, nil
}

type SimpleTextRender struct {
}

// RenderBytes implements render.
func (s *SimpleTextRender) RenderBytes(in []byte) ([]byte, error) {
	return in, nil
}

func (c *ContentView) newSimpleTextRender() (render, error) {
	return &SimpleTextRender{}, nil
}
