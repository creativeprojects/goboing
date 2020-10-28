package main

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

// PlayerPosition for left or right player
type PlayerPosition int

// PlayerPosition
const (
	PlayerLeft PlayerPosition = iota
	PlayerRight
)

// PlayerState defines when the player is playing, ready, or recovering from losing a ball
type PlayerState int

// PlayerState
const (
	PlayerStatePlaying PlayerState = iota
	PlayerStateRecovering
	PlayerStateReady
	PlayerWinningScore
)

// Scored indicates a win or a lost
type Scored int

// Scored
const (
	ScoredWin Scored = iota
	ScoredLost
)

// Player (human or AI)
type Player struct {
	// game     *Game
	normalDigits  [10]*ebiten.Image        // normal digits for displaying the score during the game
	winningDigits [10]*ebiten.Image        // winning digits for displaying the score right after scoring
	effects       [2]*ebiten.Image         // lost ball effects (left and right)
	op            *ebiten.DrawImageOptions // keep an instance available for image translation
	position      PlayerPosition           // left or right
	score         float64                  // player score
	timer         int                      // used to display the lost ball effect for n frames
	scored        Scored                   // Last score was a win or a lost
}

// NewPlayer creates a new player (left or right)
func NewPlayer(position PlayerPosition) *Player {
	winningPrefix := "digit1"
	if position == PlayerLeft {
		winningPrefix = "digit2"
	}
	return &Player{
		normalDigits: [10]*ebiten.Image{
			images["digit00"],
			images["digit01"],
			images["digit02"],
			images["digit03"],
			images["digit04"],
			images["digit05"],
			images["digit06"],
			images["digit07"],
			images["digit08"],
			images["digit09"],
		},
		winningDigits: [10]*ebiten.Image{
			images[winningPrefix+"0"],
			images[winningPrefix+"1"],
			images[winningPrefix+"2"],
			images[winningPrefix+"3"],
			images[winningPrefix+"4"],
			images[winningPrefix+"5"],
			images[winningPrefix+"6"],
			images[winningPrefix+"7"],
			images[winningPrefix+"8"],
			images[winningPrefix+"9"],
		},
		effects:  [2]*ebiten.Image{images["effect0"], images["effect1"]},
		op:       &ebiten.DrawImageOptions{},
		position: position,
		score:    0,
		timer:    0,
	}
}

// Reset player ready for a new game
func (p *Player) Reset() {
	p.score = 0
	p.timer = 0
}

// Update player state
func (p *Player) Update() {
	if p.timer < 0 {
		return
	}
	p.timer--
}

// Draw the player score
func (p *Player) Draw(screen *ebiten.Image) {
	digits := [2]int{int(math.Floor(p.score / 10)), int(math.Mod(p.score, 10))}
	images := p.normalDigits
	if p.timer > 0 && p.scored == ScoredWin {
		// Use the winning images to display the digits
		images = p.winningDigits
	}
	for i, digit := range digits {
		p.op.GeoM.Reset()
		p.op.GeoM.Translate(float64(255+(160*int(p.position))+(i*55)), 46)
		screen.DrawImage(images[digit], p.op)
	}
	// if the timer is set, it means we just lost a point and we display the effect
	if p.timer > 0 && p.scored == ScoredLost {
		screen.DrawImage(p.effects[p.position], nil)
	}
}

// BallWin indicates the player just won this point
func (p *Player) BallWin() {
	p.score++
	if p.score > 99 {
		p.score = 0
	}
	p.scored = ScoredWin
	// set the timer for 20 frames
	p.timer = 20
}

// BallLost indicates the player just lost this point
func (p *Player) BallLost() {
	p.scored = ScoredLost
	// set the timer for 20 frames
	p.timer = 20
}

// State describes if the player is (or was) playing, recovering after losing the ball, or ready to play again
func (p *Player) State() PlayerState {
	switch {
	case p.score == WinningScore:
		return PlayerWinningScore
	case p.timer < 0:
		return PlayerStatePlaying
	case p.timer > 0:
		return PlayerStateRecovering
	default:
		return PlayerStateReady
	}
}
