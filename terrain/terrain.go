package terrain

import (
	"bytes"
	"embed"
	"fmt"
	"image"
	_ "image/png"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

type Terrain struct {
	Seed                      int64
	TileSize                  int
	DecorationSpacing         float64
	DecorationSpacingVariance float64
	DecorationScaleX          float64
	DecorationScaleY          float64
	TilePaths                 []string
	TileImgs                  []*ebiten.Image
	DecorationPaths           []string
	DecorationImgs            []*ebiten.Image
}

func (t *Terrain) load(source embed.FS, paths []string) ([]*ebiten.Image, error) {
	imgs := make([]*ebiten.Image, len(paths))
	for i, path := range paths {
		i := i
		imgData, err := source.ReadFile(path)
		if err != nil {
			return imgs, fmt.Errorf("reading %s: %w", path, err)
		}
		img, _, err := image.Decode(bytes.NewReader(imgData))
		if err != nil {
			return imgs, fmt.Errorf("decoding %s: %w", path, err)
		}
		imgs[i] = ebiten.NewImageFromImage(img)
	}
	return imgs, nil
}

func (t *Terrain) Load(source embed.FS) error {

	imgs, err := t.load(source, t.TilePaths)
	if err != nil {
		return err
	}
	t.TileImgs = imgs

	imgs, err = t.load(source, t.DecorationPaths)
	if err != nil {
		return err
	}
	t.DecorationImgs = imgs

	return nil
}

func (t *Terrain) Draw(img *ebiten.Image, height float64) {
	w := img.Bounds().Dx()

	if t.Seed == 0 {
		t.Seed = time.Now().Unix()
	}
	rando := rand.New(rand.NewSource(t.Seed))

	decoXCount := (w / int(t.DecorationSpacing*t.DecorationScaleX))
	for i := 0; i < decoXCount; i++ {
		idx := rando.Int() % len(t.DecorationImgs)
		move := rando.Float64() * t.DecorationSpacingVariance
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(((float64(i)*t.DecorationSpacing)-(t.DecorationSpacingVariance/2))+move, height-float64(t.TileSize)*0.85)
		op.Filter = ebiten.FilterLinear
		op.GeoM.Scale(t.DecorationScaleX, t.DecorationScaleY)
		img.DrawImage(t.DecorationImgs[idx], op)
	}

	tileXCount := (w / t.TileSize) + 1
	for i := 0; i < tileXCount; i++ {
		idx := rando.Int() % len(t.TileImgs)
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(i*t.TileSize)-float64(t.TileSize/2), height)
		img.DrawImage(t.TileImgs[idx], op)
	}
}
