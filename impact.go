package main

import (
	"math"

	"github.com/hajimehoshi/ebiten"
)

const (
	StartImpactDuration = 0
	EndImpactDuration   = 10
	ImpactWidth         = 75.0
	ImpactHeight        = 75.0
)

// Impact animation
type Impact struct {
	images   [5]*ebiten.Image
	op       *ebiten.DrawImageOptions
	pos      Position
	duration float64
}

// NewImpact creates a new impact animation on centered coordinates
func NewImpact(x, y float64) *Impact {
	i := &Impact{
		images:   [5]*ebiten.Image{images["impact0"], images["impact1"], images["impact2"], images["impact3"], images["impact4"]},
		op:       &ebiten.DrawImageOptions{},
		pos:      NewPositionCentre(ImpactWidth, ImpactHeight, x, y),
		duration: StartImpactDuration,
	}
	return i
}

// Reset (and restart) an impact animation on centered coordinates
func (i *Impact) Reset(x, y float64) *Impact {
	i.pos = i.pos.MoveCentre(x, y)
	i.duration = StartImpactDuration
	return i
}

func (i *Impact) Update() {
	if i.HasExpired() {
		return
	}
	i.duration++
}

func (i *Impact) Draw(screen *ebiten.Image) {
	if i.HasExpired() {
		return
	}
	// We change the image every 2 ticks
	index := int(math.Floor(i.duration / 2))
	i.op.GeoM.Reset()
	i.op.GeoM.Translate(i.pos.AbsoluteX(), i.pos.AbsoluteY())
	screen.DrawImage(i.images[index], i.op)
}

// HasExpired returns true when the animation is finished
func (i *Impact) HasExpired() bool {
	return i.duration >= EndImpactDuration
}
