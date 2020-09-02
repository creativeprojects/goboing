package main

import (
	"math"
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
	game          *Game
	image         *ebiten.Image
	op            *ebiten.DrawImageOptions
	hitSounds     [5][]byte
	slowSound     []byte
	mediumSound   []byte
	fastSound     []byte
	veryfastSound []byte
	bounceSounds  [5][]byte
	borderSound   []byte
	pos           Position
	dx            float64
	dy            float64
	speed         int
}

// NewBall creates a new ball in the centre of the screen
func NewBall() *Ball {
	return (&Ball{
		image:         images["ball"],
		op:            &ebiten.DrawImageOptions{},
		hitSounds:     [5][]byte{sounds["hit0"], sounds["hit1"], sounds["hit2"], sounds["hit3"], sounds["hit4"]},
		slowSound:     sounds["hit_slow"],
		mediumSound:   sounds["hit_medium"],
		fastSound:     sounds["hit_fast"],
		veryfastSound: sounds["hit_veryfast"],
		bounceSounds:  [5][]byte{sounds["bounce0"], sounds["bounce1"], sounds["bounce2"], sounds["bounce3"], sounds["bounce4"]},
		borderSound:   sounds["bounce_synth"],
		pos:           NewPositionCentre(BallWidth, BallHeight, HalfWidth, HalfHeight),
	})
}

// Reset direction, speed and place the ball in the middle of the screen
// direction is: 1 for going right, -1 for going left
func (b *Ball) Reset(direction float64) *Ball {
	b.pos = NewPositionCentre(BallWidth, BallHeight, HalfWidth, HalfHeight)
	b.dx = direction
	b.dy = 0
	b.speed = BallStartingSpeed
	return b
}

func (b *Ball) Update() {
	// We loop to add the same increment on the ball for n times the speed
	// The collision detection runs on each incremental step so the ball is not going too far
	for n := 0; n < b.speed; n++ {
		previousX := b.pos.AbsoluteX()
		b.pos = b.pos.MoveRelative(b.dx, b.dy)
		if b.isCloseToLeftBat(previousX) {
			difference := b.pos.CentreY() - b.game.bats[PlayerLeft].CentreY()
			if b.isHittingBat(difference) {
				b.game.bats[PlayerLeft].Glow()
				b.impactAnimation()
				b.bounceFromBat(difference)
				b.playHittingBat()
				b.increaseSpeed()
				b.game.aiOffset = float64(rand.Intn(20) - 10)
				continue
			}
		}
		if b.isCloseToRightBat(previousX) {
			difference := b.pos.CentreY() - b.game.bats[PlayerRight].CentreY()
			if b.isHittingBat(difference) {
				b.game.bats[PlayerRight].Glow()
				b.impactAnimation()
				b.bounceFromBat(difference)
				b.playHittingBat()
				b.increaseSpeed()
				b.game.aiOffset = float64(rand.Intn(20) - 10)
				continue
			}
		}
		if math.Abs(b.pos.CentreY()-HalfHeight) > 220 {
			// move to the other direction
			b.dy = -b.dy
			// and get the ball out of the border
			b.pos = b.pos.MoveRelative(0, b.dy)
			b.impactAnimation()
			b.playHittingBorder()
		}
	}
}

// IsOut is true when the ball went out of the screen
func (b *Ball) IsOut() bool {
	return (b.pos.AbsoluteX() < 0.0) || (b.pos.AbsoluteX()+BallWidth > WindowWidth)
}

func (b *Ball) Draw(screen *ebiten.Image) {
	b.op.GeoM.Reset()
	b.op.GeoM.Translate(b.pos.AbsoluteX(), b.pos.AbsoluteY())
	screen.DrawImage(b.image, b.op)
}

// CentreX returns the position (on the X axis) of the centre of the ball
func (b *Ball) CentreX() float64 {
	return b.pos.CentreX()
}

// CentreY returns the position (on the Y axis) of the centre of the ball
func (b *Ball) CentreY() float64 {
	return b.pos.CentreY()
}

func (b *Ball) isCloseToLeftBat(previousX float64) bool {
	return b.pos.AbsoluteX() <= BatLeftEdge && previousX > BatLeftEdge
}

func (b *Ball) isCloseToRightBat(previousX float64) bool {
	return b.pos.AbsoluteX() >= BatRightEdge && previousX < BatRightEdge
}

func (b *Ball) playHittingBat() {
	b.game.SoundEffect(b.hitSounds[rand.Intn(4)])
	switch {
	case b.speed <= 8:
		b.game.SoundEffect(b.slowSound)
	case b.speed <= 12:
		b.game.SoundEffect(b.mediumSound)
	case b.speed <= 16:
		b.game.SoundEffect(b.fastSound)
	default:
		b.game.SoundEffect(b.veryfastSound)
	}
}

func (b *Ball) playHittingBorder() {
	b.game.SoundEffect(b.bounceSounds[rand.Intn(4)])
	b.game.SoundEffect(b.borderSound)
}

func (b *Ball) impactAnimation() {
	b.game.NewImpact(b.CentreX(), b.CentreY())
}

func (b *Ball) increaseSpeed() {
	if b.speed < BallMaxSpeed {
		b.speed++
	}
}

func (b *Ball) isHittingBat(difference float64) bool {
	// return true
	return difference > -64 && difference < 64
}

func (b *Ball) bounceFromBat(difference float64) {
	// bounce the opposite way
	b.dx = -b.dx
	// defect slightly depending on where the ball hit the bat
	b.dy += difference / 128
	// limit the Y component of the vector so we don't get into a situation where the ball is bouncing
	// up and down too rapidly
	b.dy = math.Min(math.Max(b.dy, -1), 1)
	// keep a constant speed no matter the angle
	b.dx, b.dy = Normalise(b.dx, b.dy)
}

// Normalise returns a vector with constant speed per cycle
func Normalise(x, y float64) (float64, float64) {
	length := math.Hypot(x, y)
	return x / length, y / length
}
