package main

import (
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/audio"
)

var (
	images     map[string]*ebiten.Image
	sounds     map[string][]byte
	imageNames = []string{"table", "menu0", "menu1", "over", "bat00", "bat01", "bat02", "bat10", "bat11", "bat12", "ball"}
	soundNames = []string{"down", "up", "hit0", "hit1", "hit2", "hit3", "hit4"}
	// TODO: remove this global
	game *Game
)

func main() {
	var err error
	images, err = loadImages(imageNames)
	if err != nil {
		log.Fatal(err)
	}

	audioContext, err := audio.NewContext(SampleRate)
	if err != nil {
		log.Fatal(err)
	}

	sounds, err = loadSounds(audioContext, soundNames)
	if err != nil {
		log.Fatal(err)
	}

	ebiten.SetRunnableOnUnfocused(true)
	ebiten.SetWindowSize(WindowWidth, WindowHeight)
	ebiten.SetWindowTitle(WindowTitle)
	game, err = NewGame(audioContext)
	if err != nil {
		log.Fatal(err)
	}
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
