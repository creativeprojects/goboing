package main

type PlayerType int

const (
	PlayerLeft PlayerType = iota
	PlayerRight
)

type Player struct {
	which PlayerType
	score int
}

func NewPlayer(which PlayerType) *Player {
	return &Player{
		which: which,
		score: 0,
	}
}
