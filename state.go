package main

// State is menu / playing / game over
type State int

// Current state
const (
	StateMenu State = iota
	StatePlaying
	StateGameOver
)
