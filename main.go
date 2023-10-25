package main

import (
	"embed"
	_ "image/png"
	"log"

	"github.com/DemmyDemon/go-fishing/game"
	"github.com/hajimehoshi/ebiten/v2"
)

//go:embed img/*
var data embed.FS

func main() {
	ebiten.SetWindowSize(game.ScreenWidth, game.ScreenHeight)
	ebiten.SetWindowTitle("Go Fishing")

	fishingGame, err := game.New(data)
	if err != nil {
		log.Fatal(err)
	}

	if err := ebiten.RunGame(fishingGame); err != nil {
		log.Fatal(err)
	}
}
