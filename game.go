package main

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/audio"
	"github.com/hajimehoshi/ebiten/inpututil"
)

type Game struct {
	state         State
	totalPlayers  int
	musicPlayer   *Player
	musicPlayerCh chan *Player
	errCh         chan error
}

func NewGame() (*Game, error) {
	audioContext, err := audio.NewContext(SampleRate)
	if err != nil {
		return nil, err
	}

	m, err := NewPlayer(audioContext)
	if err != nil {
		return nil, err
	}

	return &Game{
		state:         StateMenu,
		totalPlayers:  1,
		musicPlayer:   m,
		musicPlayerCh: make(chan *Player),
		errCh:         make(chan error),
	}, nil
}

func (g *Game) Update(screen *ebiten.Image) error {
	if g.state == StateMenu {
		if inpututil.IsKeyJustPressed(ebiten.KeyUp) || inpututil.IsKeyJustPressed(ebiten.KeyDown) {
			if g.totalPlayers == 1 {
				g.totalPlayers = 2
			} else {
				g.totalPlayers = 1
			}
		}

		if inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
			g.state = StatePlaying
		}
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	switch g.state {
	case StateMenu:
		if g.totalPlayers == 1 {
			screen.DrawImage(imageMenu1, nil)
		} else {
			screen.DrawImage(imageMenu2, nil)
		}
	case StatePlaying:
		screen.DrawImage(imageTable, nil)
	case StateGameOver:
		screen.DrawImage(imageGameOver, nil)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return WindowWidth, WindowHeight
}
