package water

import (
	"image"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

func New(height int, screenWidth int, screenHeight int) Water {
	water := Water{
		PointCount:   8,
		Height:       height,
		ScreenWidth:  float32(screenWidth),
		ScreenHeight: float32(screenHeight),

		// Color:        color.RGBA{78, 200, 209, 128},
		R: 78.0 / 255.0,
		G: 200.0 / 255.0,
		B: 209.0 / 255.0,
		A: 72.0 / 255.0,

		DrawOptions: &ebiten.DrawTrianglesOptions{},

		Turbulence: 5,
	}
	water.DrawOptions.AntiAlias = true
	water.DrawOptions.FillRule = ebiten.EvenOdd
	return water
}

var (
	whiteImage    = ebiten.NewImage(3, 3)
	whiteSubImage = whiteImage.SubImage(image.Rect(1, 1, 2, 2)).(*ebiten.Image)
)

func init() {
	whiteImage.Fill(color.White)
}

type Water struct {
	PointCount   int
	Height       int
	ScreenWidth  float32
	ScreenHeight float32
	Turbulence   float64
	R            float32
	G            float32
	B            float32
	A            float32
	DrawOptions  *ebiten.DrawTrianglesOptions
}

func (w *Water) maxCounter(index int) int {
	// This is entirely cargo cult. No idea what these magic numbers are.
	return 128 + (17*index+32)%64
}

func (w *Water) IndexToPoint(i int, counter int) (float32, float32) {
	x, y := float32(i)*w.ScreenWidth/float32(w.PointCount-1), float32(w.ScreenHeight/2)
	y += float32(w.Turbulence * math.Sin(float64(counter)*2*math.Pi/float64(w.maxCounter(i))))
	return x, y
}

func (w *Water) Draw(screen *ebiten.Image, counter int) {
	var path vector.Path
	for i := 0; i <= w.PointCount; i++ {
		if i == 0 {
			path.MoveTo(w.IndexToPoint(i, counter))
			continue
		}
		cpx0, cpy0 := w.IndexToPoint(i-1, counter)
		x, y := w.IndexToPoint(i, counter)
		cpx1, cpy1 := x, y
		cpx0 += 30
		cpx1 -= 30
		path.CubicTo(cpx0, cpy0, cpx1, cpy1, x, y)
	}
	path.LineTo(w.ScreenWidth, w.ScreenHeight)
	path.LineTo(0, w.ScreenHeight)

	vertices, indices := path.AppendVerticesAndIndicesForFilling(nil, nil)

	for idx := range vertices {
		vertices[idx].SrcX = 1
		vertices[idx].SrcY = 1
		vertices[idx].ColorR = w.R
		vertices[idx].ColorG = w.G
		vertices[idx].ColorB = w.B
		vertices[idx].ColorA = w.A
	}

	screen.DrawTriangles(vertices, indices, whiteSubImage, w.DrawOptions)
}
