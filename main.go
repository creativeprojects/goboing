package main

import (
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

var (
	imageTable    *ebiten.Image
	imageMenu1    *ebiten.Image
	imageMenu2    *ebiten.Image
	imageGameOver *ebiten.Image
)

func init() {
	var err error
	imageTable, _, err = ebitenutil.NewImageFromFile("images/table.png", ebiten.FilterDefault)
	if err != nil {
		log.Fatal(err)
	}
	imageMenu1, _, err = ebitenutil.NewImageFromFile("images/menu0.png", ebiten.FilterDefault)
	if err != nil {
		log.Fatal(err)
	}
	imageMenu2, _, err = ebitenutil.NewImageFromFile("images/menu1.png", ebiten.FilterDefault)
	if err != nil {
		log.Fatal(err)
	}
	imageGameOver, _, err = ebitenutil.NewImageFromFile("images/over.png", ebiten.FilterDefault)
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
