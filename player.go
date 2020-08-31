package main

import (
	"io/ioutil"
	"time"

	"github.com/hajimehoshi/ebiten/audio"
	"github.com/hajimehoshi/ebiten/audio/vorbis"
)

// Player represents the current audio state.
type Player struct {
	audioContext *audio.Context
	audioPlayer  *audio.Player
	current      time.Duration
	total        time.Duration
	seBytes      []byte
	seCh         chan []byte
	volume128    int
}

func NewPlayer(audioContext *audio.Context) (*Player, error) {
	type audioStream interface {
		audio.ReadSeekCloser
		Length() int64
	}

	const bytesPerSample = 4 // TODO: This should be defined in audio package

	var s audioStream
	var err error
	theme, err := ioutil.ReadFile("music/theme.ogg")
	if err != nil {
		return nil, err
	}
	s, err = vorbis.Decode(audioContext, audio.BytesReadSeekCloser(theme))
	if err != nil {
		return nil, err
	}

	p, err := audio.NewPlayer(audioContext, s)
	if err != nil {
		return nil, err
	}
	player := &Player{
		audioContext: audioContext,
		audioPlayer:  p,
		total:        time.Second * time.Duration(s.Length()) / bytesPerSample / SampleRate,
		volume128:    128,
		seCh:         make(chan []byte),
	}
	if player.total == 0 {
		player.total = 1
	}
	player.audioPlayer.Play()
	// go func() {
	// 	s, err := wav.Decode(audioContext, audio.BytesReadSeekCloser(raudio.Jab_wav))
	// 	if err != nil {
	// 		log.Fatal(err)
	// 		return
	// 	}
	// 	b, err := ioutil.ReadAll(s)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 		return
	// 	}
	// 	player.seCh <- b
	// }()
	return player, nil
}

func (p *Player) Close() error {
	return p.audioPlayer.Close()
}

func (p *Player) update() error {
	select {
	case p.seBytes = <-p.seCh:
		close(p.seCh)
		p.seCh = nil
	default:
	}

	if p.audioPlayer.IsPlaying() {
		p.current = p.audioPlayer.Current()
	}
	return nil
}
