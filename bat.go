package main

import (
	"math"

	"github.com/hajimehoshi/ebiten"
)

// Bat dimensions
const (
	BatWidth  = 160.0
	BatHeight = 160.0
)

var (
	batTopY    = 0.0
	batBottomY = WindowHeight - BatHeight
)

// Bat represents the left or right bat in the game
type Bat struct {
	player *Player
	images [3]*ebiten.Image
	op     *ebiten.DrawImageOptions
	pos    Position
	timer  int
	status int
	speed  float64
	moving int
}

// NewBat creates a new bat for the player
func NewBat(player *Player) *Bat {
	x := -40.0
	imagePrefix := "bat0"
	if player.position == PlayerRight {
		x = 680
		imagePrefix = "bat1"
	}
	return &Bat{
		player: player,
		images: [3]*ebiten.Image{
			images[imagePrefix+"0"],
			images[imagePrefix+"1"],
			images[imagePrefix+"2"],
		},
		op:     &ebiten.DrawImageOptions{},
		pos:    NewPositionAbsolute(BatWidth, BatHeight, x, HalfHeight-80),
		timer:  0,
		status: 0,
		moving: 0,
	}
}

// Update bat movements
func (b *Bat) Update() {
	if b.timer > 0 {
		b.timer--
	}
	// reset status back to normal
	if b.timer == 0 && b.status > 0 {
		b.status = 0
	}
}

// Draw bat on the screen
func (b *Bat) Draw(screen *ebiten.Image) {
	b.op.GeoM.Reset()
	b.op.GeoM.Translate(b.pos.AbsoluteX(), b.pos.AbsoluteY())
	screen.DrawImage(b.images[b.status], b.op)
}

// StopMoving stops the bat and reset to the initial speed
func (b *Bat) StopMoving() {
	b.moving = 0
}

// MoveUp moves the bat up at an increasing speed
func (b *Bat) MoveUp() {
	if b.moving > -1 {
		// reset the speed
		b.speed = PlayerStartSpeed
		b.moving = -1
	}
	b.moveUp(b.speed)
	b.speed++
}

func (b *Bat) moveUp(speed float64) {
	b.pos = b.pos.MoveAbsolute(b.pos.AbsoluteX(), math.Max(batTopY, b.pos.AbsoluteY()-speed))
}

// MoveDown moves the bat down at an increasing speed
func (b *Bat) MoveDown() {
	if b.moving < 1 {
		// reset the speed
		b.speed = PlayerStartSpeed
		b.moving = 1
	}
	b.moveDown(b.speed)
	b.speed++
}

func (b *Bat) moveDown(speed float64) {
	b.pos = b.pos.MoveAbsolute(b.pos.AbsoluteX(), math.Min(batBottomY, b.pos.AbsoluteY()+speed))
}

// Glow the bat when the ball just touched it
func (b *Bat) Glow() {
	b.timer = 10
	b.status = 1
}

// CentreY returns the position (on the Y axis) of the centre of the bat
func (b *Bat) CentreY() float64 {
	return b.pos.CentreY()
}

// AI player move
func (b *Bat) AI(ballX, ballY, aiOffset float64) {
	distanceX := math.Abs(ballX - b.pos.CentreX())
	targetY1 := float64(HalfHeight)
	targetY2 := ballY + aiOffset
	weight1 := math.Min(1, distanceX/HalfWidth)
	weight2 := 1 - weight1
	targetY := (weight1 * targetY1) + (weight2 * targetY2)
	move := math.Min(AIMaxSpeed, math.Max(-AIMaxSpeed, targetY-b.pos.CentreY()))
	switch {
	case move > 0:
		b.moveDown(move)
	case move < 0:
		b.moveUp(-move)
	}
}
