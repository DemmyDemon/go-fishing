package clouds

import (
	"bytes"
	"embed"
	"fmt"
	"image"
	_ "image/png"
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

func (c *Clouds) Load(source embed.FS, path string) error {
	c.CloudImages = make([]*ebiten.Image, 0)
	dir, err := source.ReadDir(path)
	if err != nil {
		return fmt.Errorf("loading clouds: %w", err)
	}
	for _, entry := range dir {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		imgData, err := source.ReadFile(path + "/" + name)
		if err != nil {
			return fmt.Errorf("reading %s/%s: %w", path, name, err)
		}
		img, _, err := image.Decode(bytes.NewReader(imgData))
		if err != nil {
			return fmt.Errorf("decoding %s/%s: %w", path, name, err)
		}
		c.CloudImages = append(c.CloudImages, ebiten.NewImageFromImage(img))
	}
	c.Clouds = make([]Cloud, c.Count)
	for i := 0; i < c.Count; i++ {
		cloud := Cloud{
			ImageIndex: rand.Intn(len(c.CloudImages)),
			Position:   rand.Float64() * 1000,
			Height:     rand.Float64() * 100,
		}
		c.Clouds[i] = cloud
	}
	return nil
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
