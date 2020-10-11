package main

// Game defaults
const (
	WindowWidth       = 800.0
	WindowHeight      = 480.0
	HalfWidth         = WindowWidth / 2
	HalfHeight        = WindowHeight / 2
	WindowTitle       = "Boing!"
	SampleRate        = 44100
	PlayerStartSpeed  = 4.0
	PlayerMaxSpeed    = 12.0
	AIMaxSpeed        = 6.0
	BallStartingSpeed = 5
	BallMaxSpeed      = 20
	WinningScore      = 10.0
	GameFullSpeed     = 60
	GameSlowSpeed     = 30
)

var (
	Player1Choice = [2][2]int{{200, 290}, {600, 340}}
	Player2Choice = [2][2]int{{200, 370}, {600, 420}}
)
