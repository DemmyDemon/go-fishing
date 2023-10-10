package main

import (
	"embed"
	"image/color"
	_ "image/png"
	"log"
	"os"
	"time"

	"github.com/DemmyDemon/go-fishing/clouds"
	"github.com/DemmyDemon/go-fishing/terrain"
	"github.com/DemmyDemon/go-fishing/water"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	screenWidth  = 1440
	screenHeight = 900
)

//go:embed img/*
var img embed.FS

type Game struct {
	time              time.Time
	delta             time.Duration
	counter           int
	backgroundHeight  int
	background        terrain.Terrain
	backgroundImage   *ebiten.Image
	backgroundOptions *ebiten.DrawImageOptions

	foregroundHeight  int
	foreground        terrain.Terrain
	foregroundImage   *ebiten.Image
	foregroundOptions *ebiten.DrawImageOptions

	water water.Water

	clouds        clouds.Clouds
	cloudsImage   *ebiten.Image
	cloudsOptions *ebiten.DrawImageOptions
}

func (g *Game) Update() error {
	now := time.Now()
	g.delta = now.Sub(g.time)
	g.time = now

	g.counter++
	/*
		if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) {
			g.clouds.Clouds[0].Drift += 10
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) {
			g.clouds.Clouds[0].Drift -= 10
		}
	*/
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		os.Exit(0)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyF5) {
		g.background.Seed = time.Now().Unix()
		g.backgroundImage = nil
		g.foregroundImage = nil
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {

	screen.Fill(color.RGBA{135, 206, 235, 255})

	if g.cloudsImage == nil {
		sky := g.clouds.Draw(screenWidth, screenHeight)
		g.cloudsImage = sky
	}
	screen.DrawImage(g.cloudsImage, g.cloudsOptions)

	if g.backgroundImage == nil {
		bg := ebiten.NewImage(screenWidth+g.background.TileSize, screenHeight)
		g.background.Draw(bg, float64(g.backgroundHeight))
		g.backgroundImage = bg
	}
	screen.DrawImage(g.backgroundImage, g.backgroundOptions)

	g.water.Draw(screen, g.counter)

	if g.foregroundImage == nil {
		g.foreground.Seed = g.background.Seed + 1
		fg := ebiten.NewImage(screenWidth+g.foreground.TileSize, screenHeight)
		g.foreground.Draw(fg, float64(g.foregroundHeight))
		g.foregroundImage = fg
	}
	screen.DrawImage(g.foregroundImage, g.foregroundOptions)

	g.water.Draw(screen, g.counter+60)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Go Fishing")

	game := &Game{
		foreground:        terrain.Sand,
		foregroundHeight:  screenHeight - (terrain.Sand.TileSize - terrain.Sand.TileSize/2),
		foregroundOptions: &ebiten.DrawImageOptions{},

		background:        terrain.Sand,
		backgroundHeight:  screenHeight - terrain.Sand.TileSize,
		backgroundOptions: &ebiten.DrawImageOptions{},

		water:         water.New(screenHeight-300, screenWidth, screenHeight),
		clouds:        clouds.Clouds{Count: 5},
		cloudsOptions: &ebiten.DrawImageOptions{},
	}
	game.time = time.Now()

	game.backgroundOptions.ColorScale.ScaleAlpha(0.75)
	game.backgroundOptions.GeoM.Translate(-float64(game.background.TileSize)/2, 0)
	game.foregroundOptions.GeoM.Translate(-float64(game.foreground.TileSize)/2, 0)
	game.cloudsOptions.ColorScale.ScaleAlpha(0.5)

	err := game.foreground.Load(img)
	if err != nil {
		log.Fatal(err)
	}
	err = game.background.Load(img)
	if err != nil {
		log.Fatal(err)
	}
	err = game.clouds.Load(img, "img/clouds")
	if err != nil {
		log.Fatal(err)
	}
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
