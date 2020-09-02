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

const (
	StartImpacts = 5
)

var (
	playersSelection = 1
)

// Game contains the current game state
type Game struct {
	audioContext *audio.Context
	musicPlayer  *AudioPlayer
	state        GameState
	totalPlayers int
	players      [2]*Player
	bats         [2]*Bat
	ball         *Ball
	impacts      []*Impact
	debug        bool
	aiOffset     float64
}

// NewGame creates a new game instance and prepares a demo AI game
func NewGame(audioContext *audio.Context) (*Game, error) {

	m, err := NewAudioPlayer(audioContext)
	if err != nil {
		return nil, err
	}
	ball := NewBall()
	g := &Game{
		audioContext: audioContext,
		musicPlayer:  m,
		state:        StateMenu,
		bats: [2]*Bat{
			NewBat(NewPlayer(PlayerLeft)),
			NewBat(NewPlayer(PlayerRight)),
		},
		ball:    ball,
		impacts: make([]*Impact, 0, StartImpacts),
	}
	// circular references
	ball.game = g

	direction := math.Round(rand.Float64())
	if direction == 0 {
		direction = -1
	}
	g.ball.Reset(direction)

	return g, nil
}

// Layout defines the size of the game in pixels
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return WindowWidth, WindowHeight
}

// Start initializes a game with a number of players
func (g *Game) Start(players int) *Game {
	if players < 0 || players > 2 {
		players = 0
	}
	g.totalPlayers = players

	direction := math.Round(rand.Float64())
	if direction == 0 {
		direction = -1
	}
	g.ball.Reset(direction)
	g.state = StatePlaying
	return g
}

// Update game events
func (g *Game) Update(screen *ebiten.Image) error {
	// Debug
	if inpututil.IsKeyJustPressed(ebiten.KeyD) {
		g.debug = !g.debug
	}

	if g.state == StateMenu {

		if inpututil.IsKeyJustPressed(ebiten.KeyDown) && playersSelection == 1 {
			playersSelection = 2
			PlaySE(g.audioContext, sounds["down"])
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyUp) && playersSelection == 2 {
			playersSelection = 1
			PlaySE(g.audioContext, sounds["up"])
		}

		if inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
			g.Start(playersSelection)
		}

		// Escape
		if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
			g.Start(0)
		}
		return nil
	}
	if g.state == StatePlaying {
		// Escape
		if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
			g.state = StateMenu
		}
		// Pause
		if inpututil.IsKeyJustPressed(ebiten.KeyP) {
			g.state = StatePaused
		}

		if g.totalPlayers > 0 {
			if ebiten.IsKeyPressed(ebiten.KeyA) {
				g.bats[PlayerLeft].MoveUp(PlayerSpeed)
			}
			if ebiten.IsKeyPressed(ebiten.KeyZ) {
				g.bats[PlayerLeft].MoveDown(PlayerSpeed)
			}
		} else {
			g.bats[PlayerLeft].AI(g.ball.pos.CentreX(), g.ball.pos.CentreY(), g.aiOffset)
		}
		if g.totalPlayers > 1 {
			if ebiten.IsKeyPressed(ebiten.KeyK) {
				g.bats[PlayerRight].MoveUp(PlayerSpeed)
			}
			if ebiten.IsKeyPressed(ebiten.KeyM) {
				g.bats[PlayerRight].MoveDown(PlayerSpeed)
			}
		} else {
			g.bats[PlayerRight].AI(g.ball.pos.CentreX(), g.ball.pos.CentreY(), g.aiOffset)
		}
		// run impacts first
		for _, impact := range g.impacts {
			impact.Update()
		}

		for _, bat := range g.bats {
			bat.Update()
		}
		g.ball.Update()

		if g.ball.IsOut() {
			g.SoundEffect(sounds["score_goal"])
			g.ball.Reset(1)
		}
		return nil
	}

	if g.state == StatePaused {
		if inpututil.IsKeyJustPressed(ebiten.KeyP) {
			g.state = StatePlaying
		}
		return nil
	}
	return nil
}

// Draw game events
func (g *Game) Draw(screen *ebiten.Image) {
	screen.DrawImage(images["table"], nil)

	switch g.state {
	case StateMenu:
		if playersSelection == 2 {
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

	for _, impact := range g.impacts {
		impact.Draw(screen)
	}

	if g.debug {
		g.displayDebug(screen)
	}
}

// SoundEffect plays a sound in the game
func (g *Game) SoundEffect(se []byte) {
	PlaySE(g.audioContext, se)
}

// NewImpact adds an impact animation at the coordinates (use centered coordinates)
func (g *Game) NewImpact(x, y float64) {
	// Reuse a free impact first
	for index := 0; index < len(g.impacts); index++ {
		if g.impacts[index].HasExpired() {
			g.impacts[index].Reset(x, y)
			return
		}
	}
	// No one was available
	g.impacts = append(g.impacts, NewImpact(x, y))
}

func (g *Game) displayDebug(screen *ebiten.Image) {
	template := " TPS: %0.2f \n Left bat: %0.0f \n Right bat: %0.0f \n Ball: %0.0f, %0.0f \n Impacts: %d \n"
	msg := fmt.Sprintf(template,
		ebiten.CurrentTPS(),
		g.bats[0].CentreY(),
		g.bats[1].CentreY(),
		g.ball.CentreX(),
		g.ball.CentreY(),
		len(g.impacts),
	)
	ebitenutil.DebugPrint(screen, msg)
}
