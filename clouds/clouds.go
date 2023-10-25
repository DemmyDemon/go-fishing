package clouds

import (
	"fmt"
	"image"
	_ "image/png"
	"io/fs"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
)

type Cloud struct {
	ImageIndex int
	Position   float64
	Height     float64
}

type Clouds struct {
	Count       int
	CloudImages []*ebiten.Image
	Clouds      []Cloud
}

func (c *Clouds) Load(source fs.FS, path string) error {
	c.CloudImages = make([]*ebiten.Image, 0)
	err := fs.WalkDir(source, path, func(path string, entry fs.DirEntry, fileError error) error {
		if entry.IsDir() {
			return nil
		}
		if fileError != nil {
			return fileError
		}

		name := entry.Name()
		imgData, fileError := source.Open(path)
		if fileError != nil {
			return fmt.Errorf("reading %s/%s: %w", path, name, fileError)
		}
		img, _, fileError := image.Decode(imgData)
		if fileError != nil {
			return fmt.Errorf("decoding %s/%s: %w", path, name, fileError)
		}
		c.CloudImages = append(c.CloudImages, ebiten.NewImageFromImage(img))
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (c *Clouds) RandomizePositions() {
	c.Clouds = make([]Cloud, c.Count)
	for i := 0; i < c.Count; i++ {
		cloud := Cloud{
			ImageIndex: rand.Intn(len(c.CloudImages)),
			Position:   rand.Float64() * 1000,
			Height:     rand.Float64() * 100,
		}
		c.Clouds[i] = cloud
	}
}

func (c *Cloud) Draw(screen *ebiten.Image, cloudImage *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(0.5, 0.5)
	op.GeoM.Translate(c.Position, c.Height)
	screen.DrawImage(cloudImage, op)
}

func (c *Clouds) Draw(screenWidth int, screenHeight int) *ebiten.Image {
	sky := ebiten.NewImage(screenWidth, screenHeight)
	for _, cloud := range c.Clouds {
		cloud.Draw(sky, c.CloudImages[cloud.ImageIndex])
	}
	return sky
}
