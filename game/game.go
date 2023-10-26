package game

import (
	"image/color"
	"io/fs"
	"os"
	"time"

	"github.com/DemmyDemon/go-fishing/clouds"
	"github.com/DemmyDemon/go-fishing/hook"
	"github.com/DemmyDemon/go-fishing/terrain"
	"github.com/DemmyDemon/go-fishing/water"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	ScreenWidth  = 1600
	ScreenHeight = 900
)

type Game struct {
	Time              time.Time
	delta             time.Duration
	counter           int
	BackgroundHeight  int
	Background        terrain.Terrain
	BackgroundImage   *ebiten.Image
	BackgroundOptions *ebiten.DrawImageOptions

	ForegroundHeight  int
	Foreground        terrain.Terrain
	ForegroundImage   *ebiten.Image
	ForegroundOptions *ebiten.DrawImageOptions

	Water water.Water

	Clouds        clouds.Clouds
	CloudsImage   *ebiten.Image
	CloudsOptions *ebiten.DrawImageOptions

	Hook hook.Hook
}

func New(data fs.FS) (*Game, error) {
	game := &Game{
		Foreground:        terrain.Sand,
		ForegroundHeight:  ScreenHeight - (terrain.Sand.TileSize - terrain.Sand.TileSize/2),
		ForegroundOptions: &ebiten.DrawImageOptions{},

		Background:        terrain.Sand,
		BackgroundHeight:  ScreenHeight - terrain.Sand.TileSize,
		BackgroundOptions: &ebiten.DrawImageOptions{},

		Water:         water.New(ScreenHeight-300, ScreenWidth, ScreenHeight),
		Clouds:        clouds.Clouds{Count: 8},
		CloudsOptions: &ebiten.DrawImageOptions{},

		Hook: hook.New(),
	}
	game.Time = time.Now()

	game.BackgroundOptions.ColorScale.ScaleAlpha(0.5)
	game.BackgroundOptions.GeoM.Translate(-float64(game.Background.TileSize)/2, 0)
	game.ForegroundOptions.GeoM.Translate(-float64(game.Foreground.TileSize)/2, 0)
	game.CloudsOptions.ColorScale.ScaleAlpha(0.75)

	err := game.Foreground.Load(data)
	if err != nil {
		return nil, err
	}
	err = game.Background.Load(data)
	if err != nil {
		return nil, err
	}
	err = game.Clouds.Load(data, "img/clouds")
	if err != nil {
		return nil, err
	}
	err = game.Hook.Load(data, "img/hook/hook.png")
	if err != nil {
		return nil, err
	}

	return game, nil
}

func (g *Game) Update() error {
	now := time.Now()
	g.delta = now.Sub(g.Time)
	g.Time = now

	g.counter++

	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) || inpututil.IsKeyJustPressed(ebiten.KeyQ) {
		os.Exit(0)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyF5) {
		g.Background.Seed = time.Now().Unix()
		g.BackgroundImage = nil
		g.ForegroundImage = nil
		g.CloudsImage = nil
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyF) {
		ebiten.SetFullscreen(!ebiten.IsFullscreen())
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {

	screen.Fill(color.RGBA{135, 206, 235, 255})

	if g.CloudsImage == nil {
		g.Clouds.RandomizePositions()
		sky := g.Clouds.Draw(ScreenWidth, ScreenHeight)
		g.CloudsImage = sky
	}
	screen.DrawImage(g.CloudsImage, g.CloudsOptions)

	if g.BackgroundImage == nil {
		bg := ebiten.NewImage(ScreenWidth+g.Background.TileSize, ScreenHeight)
		g.Background.Draw(bg, float64(g.BackgroundHeight))
		g.BackgroundImage = bg
	}
	screen.DrawImage(g.BackgroundImage, g.BackgroundOptions)

	g.Water.Draw(screen, g.counter)

	if g.ForegroundImage == nil {
		g.Foreground.Seed = g.Background.Seed + 1
		fg := ebiten.NewImage(ScreenWidth+g.Foreground.TileSize, ScreenHeight)
		g.Foreground.Draw(fg, float64(g.ForegroundHeight))
		g.ForegroundImage = fg
	}
	screen.DrawImage(g.ForegroundImage, g.ForegroundOptions)

	g.Hook.Draw(screen, 500, 200)

	g.Water.Draw(screen, g.counter+60)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return ScreenWidth, ScreenHeight
}
