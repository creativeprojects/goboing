package main

import (
	"fmt"
	"math"
	"math/rand"

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
	bats          [2]*Bat
	ball          *Ball
}

func NewGame(audioContext *audio.Context) (*Game, error) {

	m, err := NewAudioPlayer(audioContext)
	if err != nil {
		return nil, err
	}

	direction := math.Round(rand.Float64())
	if direction == 0 {
		direction = -1
	}
	return &Game{
		audioContext:  audioContext,
		state:         StateMenu,
		totalPlayers:  1,
		musicPlayer:   m,
		musicPlayerCh: make(chan *AudioPlayer),
		errCh:         make(chan error),
		bats: [2]*Bat{
			NewBat(NewPlayer(PlayerLeft)),
			NewBat(NewPlayer(PlayerRight)),
		},
		ball: NewBall(direction),
	}, nil
}

func (g *Game) Update(screen *ebiten.Image) error {
	if g.state == StateMenu {
		// if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		// 	ebiten.
		// }

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
	if g.state == StatePlaying {
		if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
			g.state = StateMenu
		}

		if ebiten.IsKeyPressed(ebiten.KeyA) {
			g.bats[0].MoveUp(PlayerSpeed)
		}
		if ebiten.IsKeyPressed(ebiten.KeyZ) {
			g.bats[0].MoveDown(PlayerSpeed)
		}
		if ebiten.IsKeyPressed(ebiten.KeyK) {
			g.bats[1].MoveUp(PlayerSpeed)
		}
		if ebiten.IsKeyPressed(ebiten.KeyM) {
			g.bats[1].MoveDown(PlayerSpeed)
		}

		// Restart the ball
		if inpututil.IsKeyJustPressed(ebiten.KeyR) {
			direction := math.Round(rand.Float64())
			if direction == 0 {
				direction = -1
			}
			g.ball = NewBall(direction)
		}

		if g.ball.IsOut() {
			return nil
		}
		g.ball.Update()
		for _, bat := range g.bats {
			bat.Update()
		}
	}

	if g.state == StatePaused {
		if inpututil.IsKeyJustPressed(ebiten.KeyP) {
			g.state = StatePlaying
		}
		// Restart the ball
		if inpututil.IsKeyJustPressed(ebiten.KeyR) {
			direction := math.Round(rand.Float64())
			if direction == 0 {
				direction = -1
			}
			g.ball = NewBall(direction)
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

	for _, bat := range g.bats {
		bat.Draw(screen)
	}

	g.ball.Draw(screen)

	msg := fmt.Sprintf(`TPS: %0.2f`, ebiten.CurrentTPS())
	ebitenutil.DebugPrint(screen, msg)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return WindowWidth, WindowHeight
}
