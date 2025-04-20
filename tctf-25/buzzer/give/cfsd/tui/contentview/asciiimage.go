package contentview

import (
	"bytes"
	"image"
	"math"

	_ "image/jpeg"
	_ "image/png"

	"github.com/qeesung/image2ascii/convert"
)

type ImageRender struct {
	converter *convert.ImageConverter

	height int
	width  int
}

// RenderBytes implements render.
func (i *ImageRender) RenderBytes(in []byte) ([]byte, error) {

	img, _, err := image.Decode(bytes.NewReader(in))
	if err != nil {
		return nil, err
	}
	imgWidth := float64(img.Bounds().Max.X)
	imgHeight := float64(img.Bounds().Max.Y)

	opts := convert.DefaultOptions

	// always fit by width
	opts.FixedWidth = i.width

	// safe proportions
	realHeight := (float64(i.width) / imgWidth) * imgHeight

	terminalRatio := 28.0 / 15.0 // NB: letters in the terminal are not square
	opts.FixedHeight = int(math.Floor(realHeight / terminalRatio))

	opts.Colored = true

	asciiImg := i.converter.Image2ASCIIString(img, &opts)
	return []byte(asciiImg), nil
}

func (c *ContentView) newImageRender() (render, error) {
	return &ImageRender{
		converter: convert.NewImageConverter(),
		height:    c.view.Height,
		width:     c.view.Width,
	}, nil
}
