package main

import (
	"fmt"
	"math/rand"

	"github.com/hajimehoshi/ebiten"
)

const (
	BallWidth    = 24.0
	BallHeight   = 24.0
	BatLeftEdge  = 42.0
	BatRightEdge = WindowWidth - 66
)

type Ball struct {
	image *ebiten.Image
	op    *ebiten.DrawImageOptions
	x     float64
	y     float64
	dx    float64
	dy    float64
	speed int
}

// NewBall creates a new ball in the center
// direction is: 1 for going right, -1 for going left
func NewBall(direction float64) *Ball {
	return &Ball{
		image: images["ball"],
		op:    &ebiten.DrawImageOptions{},
		x:     HalfWidth - 12.0,
		y:     HalfHeight - 12.0,
		dx:    direction,
		dy:    0,
		speed: 5,
	}
}

func (b *Ball) Update() {
	// We loop to add the same increment on the ball for n times the speed
	// The collision detection runs on each incremental step so the ball is not going too far
	for n := 0; n < b.speed; n++ {
		previousX := b.x
		b.x += b.dx
		b.y += b.dy
		if b.isCloseToLeftBat(previousX) {
			game.bats[0].Glow()
			b.dx = -b.dx
			PlaySE(game.audioContext, sounds[fmt.Sprintf("hit%d", rand.Intn(4))])
			continue
		}
		if b.isCloseToRightBat(previousX) {
			game.bats[1].Glow()
			b.dx = -b.dx
			PlaySE(game.audioContext, sounds[fmt.Sprintf("hit%d", rand.Intn(4))])
			continue
		}
	}
}

func (b *Ball) IsOut() bool {
	return (b.x < 0.0) || (b.x+BallWidth > WindowWidth)
}

func (b *Ball) Draw(screen *ebiten.Image) {
	b.op.GeoM.Reset()
	b.op.GeoM.Translate(b.x, b.y)
	screen.DrawImage(b.image, b.op)
}

func (b *Ball) isCloseToLeftBat(previousX float64) bool {
	return b.x <= BatLeftEdge && previousX > BatLeftEdge
}

func (b *Ball) isCloseToRightBat(previousX float64) bool {
	return b.x >= BatRightEdge && previousX < BatRightEdge
}
