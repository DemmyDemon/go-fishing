package hook

import (
	"fmt"
	"image"
	"image/color"
	_ "image/png"
	"io/fs"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Line struct {
	Color  color.Color
	Offset float32
	Width  float32
}

var AquaLine = Line{
	Color:  color.NRGBA{0, 255, 181, 191},
	Offset: 8,
	Width:  3,
}

type Hook struct {
	Image *ebiten.Image
	Scale float64
	Line  Line
}

func New() Hook {
	return Hook{
		Scale: 0.5,
		Line:  AquaLine,
	}
}

func (h *Hook) Load(source fs.FS, path string) error {
	file, err := source.Open(path)
	if err != nil {
		return fmt.Errorf("opening %s: %w", path, err)
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		return fmt.Errorf("decoding %s: %w", path, err)
	}
	h.Image = ebiten.NewImageFromImage(img)
	return nil
}

func (h *Hook) Draw(screen *ebiten.Image, x, y float32) {

	if y > 0 {
		lineX := x + h.Line.Offset
		vector.StrokeLine(screen, lineX, 0, lineX, y, h.Line.Width, h.Line.Color, false)
	}

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(0.5, 0.5)
	op.Filter = ebiten.FilterNearest
	op.GeoM.Translate(float64(x), float64(y))

	screen.DrawImage(h.Image, op)

}
