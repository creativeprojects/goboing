package main

import (
	"fmt"
	"math"
	"math/rand"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
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
	slow         bool
}

// NewGame creates a new game instance and prepares a demo AI game
func NewGame(audioContext *audio.Context) (*Game, error) {

	m, err := NewAudioPlayer(audioContext)
	if err != nil {
		return nil, err
	}
	ball := NewBall()
	player1 := NewPlayer(PlayerLeft)
	player2 := NewPlayer(PlayerRight)
	g := &Game{
		audioContext: audioContext,
		musicPlayer:  m,
		state:        StateMenu,
		players:      [2]*Player{player1, player2},
		bats:         [2]*Bat{NewBat(player1), NewBat(player2)},
		ball:         ball,
		impacts:      make([]*Impact, 0, StartImpacts),
		debug:        false,
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
func (g *Game) Update() error {
	// Debug
	if inpututil.IsKeyJustPressed(ebiten.KeyD) {
		g.debug = !g.debug
	}

	if g.state == StateMenu {
		g.UpdateMenu()
		return nil
	}
	if g.state == StatePlaying {
		g.UpdatePlaying()
		return nil
	}

	if g.state == StatePaused {
		// un-pause
		if inpututil.IsKeyJustPressed(ebiten.KeyP) || inpututil.IsKeyJustPressed(ebiten.KeySpace) {
			g.state = StatePlaying
		}
		return nil
	}

	if g.state == StateGameOver {
		// un-pause
		if inpututil.IsKeyJustPressed(ebiten.KeyEnter) || inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
			g.Reset()
			g.state = StateMenu
		}
		return nil
	}
	return nil
}

func (g *Game) UpdateMenu() {
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

	// Now check for touch events
	if ids := inpututil.JustPressedTouchIDs(); ids != nil && len(ids) > 0 {
		for _, id := range ids {
			x, y := ebiten.TouchPosition(id)
			if touchInRectangle(x, y, Player1Choice) {
				if playersSelection == 1 {
					// Start the game
					g.Start(playersSelection)
					continue
				}
				playersSelection = 1
				PlaySE(g.audioContext, sounds["up"])
				continue
			}
			if touchInRectangle(x, y, Player2Choice) {
				if playersSelection == 2 {
					// Start the game
					g.Start(playersSelection)
					continue
				}
				playersSelection = 2
				PlaySE(g.audioContext, sounds["down"])
				continue
			}
		}
	}
}

func (g *Game) UpdatePlaying() error {
	// Escape
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		g.Reset()
		g.state = StateMenu
	}
	// Pause
	if inpututil.IsKeyJustPressed(ebiten.KeyP) || inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		g.state = StatePaused
	}

	// toggle between slow and normal mode
	if inpututil.IsKeyJustPressed(ebiten.KeyS) {
		if g.slow {
			ebiten.SetMaxTPS(GameSlowSpeed)
		} else {
			ebiten.SetMaxTPS(GameFullSpeed)
		}
		g.slow = !g.slow
	}

	// get touch events if any
	hasLeftTouch, hasRightTouch := false, false
	leftTouchY, rightTouchY := 0, 0
	touchIDs := ebiten.TouchIDs()
	if touchIDs != nil && len(touchIDs) > 0 {
		for _, touchID := range touchIDs {
			// we only accept one touch on each side of the screen
			x, y := ebiten.TouchPosition(touchID)
			if x < 200 {
				if hasLeftTouch {
					continue
				}
				hasLeftTouch = true
				leftTouchY = y
			} else if x > 600 {
				if hasRightTouch {
					continue
				}
				hasRightTouch = true
				rightTouchY = y
			}
		}
	}

	if g.totalPlayers > 0 {
		if ebiten.IsKeyPressed(ebiten.KeyA) {
			g.bats[PlayerLeft].MoveUp()
		} else if ebiten.IsKeyPressed(ebiten.KeyZ) {
			g.bats[PlayerLeft].MoveDown()
		} else {
			g.bats[PlayerLeft].StopMoving()
		}

		// left touch?
		if hasLeftTouch {
			g.bats[PlayerLeft].MoveTo(float64(leftTouchY - 80))
		}
	} else {
		g.bats[PlayerLeft].AI(g.ball.sprite.X(XCentre), g.ball.sprite.Y(YCentre), g.aiOffset)
	}
	if g.totalPlayers > 1 {
		if ebiten.IsKeyPressed(ebiten.KeyK) {
			g.bats[PlayerRight].MoveUp()
		} else if ebiten.IsKeyPressed(ebiten.KeyM) {
			g.bats[PlayerRight].MoveDown()
		} else {
			g.bats[PlayerRight].StopMoving()
		}

		// right touch?
		if hasRightTouch {
			g.bats[PlayerRight].MoveTo(float64(rightTouchY - 80))
		}
	} else {
		g.bats[PlayerRight].AI(g.ball.sprite.X(XCentre), g.ball.sprite.Y(YCentre), g.aiOffset)
	}
	// run impacts first
	for _, impact := range g.impacts {
		impact.Update()
	}

	for _, bat := range g.bats {
		bat.Update()
	}
	g.ball.Update()

	for _, player := range g.players {
		player.Update()
	}

	if g.ball.IsOut() {
		var scoringPlayer, losingPlayer PlayerPosition
		if g.ball.sprite.X(XLeft) < HalfWidth {
			scoringPlayer = PlayerRight
			losingPlayer = PlayerLeft
		} else {
			scoringPlayer = PlayerLeft
			losingPlayer = PlayerRight
		}
		if g.players[losingPlayer].State() == PlayerStatePlaying {
			g.players[scoringPlayer].BallWin()
			g.players[losingPlayer].BallLost()
			g.SoundEffect(sounds["score_goal"])
		}
		if g.players[losingPlayer].State() == PlayerStateReady {
			direction := 1.0
			if losingPlayer == PlayerLeft {
				direction = -1
			}
			g.ball.Reset(direction)
		}
		if g.players[scoringPlayer].State() == PlayerWinningScore {
			// Game finished!
			g.state = StateGameOver
		}
	}
	return nil
}

// Reset game ready for a new one
func (g *Game) Reset() {
	for _, player := range g.players {
		player.Reset()
	}
	g.ball.Reset(-1)
}

// Draw game events
func (g *Game) Draw(screen *ebiten.Image) {
	screen.DrawImage(images[imageTable], nil)

	switch g.state {
	case StateMenu:
		if playersSelection == 2 {
			screen.DrawImage(images[menuTwoPlayers], nil)
		} else {
			screen.DrawImage(images[menuOnePlayer], nil)
		}
	case StatePlaying:

	case StateGameOver:
		screen.DrawImage(images[gameOver], nil)
	}

	for _, bat := range g.bats {
		bat.Draw(screen)
	}

	g.ball.Draw(screen)

	for _, impact := range g.impacts {
		impact.Draw(screen)
	}

	for _, player := range g.players {
		player.Draw(screen)
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
			g.impacts[index].Start(x, y)
			return
		}
	}
	// No one was available
	g.impacts = append(g.impacts, NewImpact().Start(x, y))
}

func (g *Game) displayDebug(screen *ebiten.Image) {
	template := " TPS: %0.2f \n Left bat: %0.0f \n Right bat: %0.0f \n Ball: %0.0f, %0.0f | %0.0f, %0.0f\n Impacts: %d \n Score: %0.0f / %0.0f \n Touch: %s"
	msg := fmt.Sprintf(template,
		ebiten.CurrentTPS(),
		g.bats[0].CentreY(),
		g.bats[1].CentreY(),
		g.ball.CentreX(),
		g.ball.CentreY(),
		g.ball.sprite.X(XCentre),
		g.ball.sprite.Y(YCentre),
		len(g.impacts),
		g.players[PlayerLeft].score,
		g.players[PlayerRight].score,
		debugTouch(),
	)
	ebitenutil.DebugPrint(screen, msg)
}

func debugTouch() string {
	ids := ebiten.TouchIDs()
	if ids == nil || len(ids) == 0 {
		return "n/a"
	}
	values := make([]string, 0, len(ids))
	for _, id := range ids {
		x, y := ebiten.TouchPosition(id)
		values = append(values, fmt.Sprintf("(%d, %d)", x, y))
	}
	return strings.Join(values, " - ")
}

func touchInRectangle(x, y int, rectangle [2][2]int) bool {
	return x >= rectangle[0][0] && x <= rectangle[1][0] &&
		y >= rectangle[0][1] && y <= rectangle[1][1]
}
