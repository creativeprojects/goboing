package main

import (
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten"
)

var (
	images map[string]*ebiten.Image
)

func init() {
	var err error
	images, err = loadImages([]string{"table.png", "menu0.png", "menu1.png", "over.png"})
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	ebiten.SetWindowSize(WindowWidth, WindowHeight)
	ebiten.SetWindowTitle(WindowTitle)
	game, err := NewGame()
	if err != nil {
		log.Fatal(err)
	}
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
