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
	imageNames = []string{"table", "menu0", "menu1", "over", "bat00", "bat01", "bat02", "bat10", "bat11", "bat12", "ball", "impact0", "impact1", "impact2", "impact3", "impact4", "digit00", "digit01", "digit02", "digit03", "digit04", "digit05", "digit06", "digit07", "digit08", "digit09", "digit10", "digit11", "digit12", "digit13", "digit14", "digit15", "digit16", "digit17", "digit18", "digit19", "digit20", "digit21", "digit22", "digit23", "digit24", "digit25", "digit26", "digit27", "digit28", "digit29", "effect0", "effect1"}
	soundNames = []string{"down", "up", "hit0", "hit1", "hit2", "hit3", "hit4", "hit_slow", "hit_medium", "hit_fast", "hit_veryfast", "score_goal", "bounce0", "bounce1", "bounce2", "bounce3", "bounce4", "bounce_synth"}
)

func main() {
	var err error

	log.SetFlags(log.LstdFlags | log.Lshortfile)

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
	game, err := NewGame(audioContext)
	if err != nil {
		log.Fatal(err)
	}
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
