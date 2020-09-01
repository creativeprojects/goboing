package main

import "github.com/hajimehoshi/ebiten"

type Impact struct {
	images   [4]*ebiten.Image
	x        float64
	y        float64
	duration int
}

func (i *Impact) Update() {
	i.duration++
}
