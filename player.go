package main

// PlayerPosition for left or right player
type PlayerPosition int

const (
	PlayerLeft PlayerPosition = iota
	PlayerRight
)

type Player struct {
	game     *Game
	position PlayerPosition
	score    int
}

func NewPlayer(position PlayerPosition) *Player {
	return &Player{
		position: position,
		score:    0,
	}
}
