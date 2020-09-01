package main

import (
	"fmt"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/audio"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/hajimehoshi/ebiten/inpututil"
)

type Game struct {
	state         State
	totalPlayers  int
	audioContext  *audio.Context
	musicPlayer   *AudioPlayer
	musicPlayerCh chan *AudioPlayer
	errCh         chan error
}

func NewGame(audioContext *audio.Context) (*Game, error) {

	m, err := NewPlayer(audioContext)
	if err != nil {
		return nil, err
	}

	return &Game{
		audioContext:  audioContext,
		state:         StateMenu,
		totalPlayers:  1,
		musicPlayer:   m,
		musicPlayerCh: make(chan *AudioPlayer),
		errCh:         make(chan error),
	}, nil
}

func (g *Game) Update(screen *ebiten.Image) error {
	if g.state == StateMenu {
		if inpututil.IsKeyJustPressed(ebiten.KeyDown) && g.totalPlayers == 1 {
			g.totalPlayers = 2
			PlaySE(g.audioContext, sounds["down"])
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyUp) && g.totalPlayers == 2 {
			g.totalPlayers = 1
			PlaySE(g.audioContext, sounds["up"])
		}

		if inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
			g.state = StatePlaying
		}
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.DrawImage(images["table"], nil)

	switch g.state {
	case StateMenu:
		if g.totalPlayers == 2 {
			screen.DrawImage(images["menu1"], nil)
		} else {
			screen.DrawImage(images["menu0"], nil)
		}
	case StatePlaying:

	case StateGameOver:
		screen.DrawImage(images["over"], nil)
	}

	msg := fmt.Sprintf(`TPS: %0.2f`, ebiten.CurrentTPS())
	ebitenutil.DebugPrint(screen, msg)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return WindowWidth, WindowHeight
}
