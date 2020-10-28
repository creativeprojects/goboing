package main

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// Impact animation
type Impact struct {
	sprite *Sprite
}

// NewImpact creates a new impact animation ready to go
func NewImpact() *Impact {
	i := &Impact{
		sprite: NewSprite(XCentre, YCentre).Animation(
			[]*ebiten.Image{images["impact0"], images["impact1"], images["impact2"], images["impact3"], images["impact4"]},
			nil, 2, false),
	}
	return i
}

// Start (or restart) an impact animation on centered coordinates
func (i *Impact) Start(x, y float64) *Impact {
	i.sprite.MoveTo(x, y).Start()
	return i
}

func (i *Impact) Update() {
	if i.HasExpired() {
		return
	}
	i.sprite.Update()
}

func (i *Impact) Draw(screen *ebiten.Image) {
	if i.HasExpired() {
		return
	}
	i.sprite.Draw(screen)
}

// HasExpired returns true when the animation is finished
func (i *Impact) HasExpired() bool {
	return i.sprite.IsFinished()
}
