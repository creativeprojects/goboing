package main

import (
	"math"

	"github.com/hajimehoshi/ebiten"
)

const (
	BatWidth  = 160.0
	BatHeight = 160.0
)

var (
	batTopY    = 0.0
	batBottomY = WindowHeight - BatHeight
)

type Bat struct {
	player *Player
	images [3]*ebiten.Image
	op     *ebiten.DrawImageOptions
	x      float64
	y      float64
	timer  int
	status int
}

func NewBat(player *Player) *Bat {
	x := -40.0
	imagePrefix := "bat0"
	if player.which == PlayerRight {
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
		x:      x,
		y:      HalfHeight - 80,
		timer:  0,
		status: 0,
	}
}

func (b *Bat) Update() {
	if b.timer > 0 {
		b.timer--
	}
	// reset status back to normal
	if b.timer == 0 && b.status > 0 {
		b.status = 0
	}
}

func (b *Bat) Draw(screen *ebiten.Image) {
	b.op.GeoM.Reset()
	b.op.GeoM.Translate(b.x, b.y)
	screen.DrawImage(b.images[b.status], b.op)
}

func (b *Bat) MoveUp(speed float64) {
	b.y = math.Max(batTopY, b.y-speed)
}

func (b *Bat) MoveDown(speed float64) {
	b.y = math.Min(batBottomY, b.y+speed)
}

func (b *Bat) Glow() {
	b.timer = 10
	b.status = 1
}
