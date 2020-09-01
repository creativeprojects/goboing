package main

import (
	"io/ioutil"
	"log"
	"time"

	"github.com/hajimehoshi/ebiten/audio"
	"github.com/hajimehoshi/ebiten/audio/vorbis"
)

// AudioPlayer represents the current audio state.
type AudioPlayer struct {
	audioContext *audio.Context
	audioPlayer  *audio.Player
	current      time.Duration
	total        time.Duration
	seBytes      []byte
	seCh         chan []byte
	volume128    int
}

func NewPlayer(audioContext *audio.Context) (*AudioPlayer, error) {
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
	player := &AudioPlayer{
		audioContext: audioContext,
		audioPlayer:  p,
		total:        time.Second * time.Duration(s.Length()) / bytesPerSample / SampleRate,
		volume128:    32,
		seCh:         make(chan []byte),
	}
	if player.total == 0 {
		player.total = 1
	}
	player.audioPlayer.SetVolume(float64(player.volume128) / 128)
	player.audioPlayer.Play()

	return player, nil
}

func (p *AudioPlayer) Close() error {
	return p.audioPlayer.Close()
}

func (p *AudioPlayer) update() error {
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

// PlaySE plays a sound effect.
func PlaySE(audioContext *audio.Context, bs []byte) {
	if bs == nil || len(bs) == 0 {
		log.Printf("cannot play empty sound")
		return
	}
	sePlayer, err := audio.NewPlayerFromBytes(audioContext, bs)
	if err != nil {
		log.Printf("error playing sound effect: %v", err)
	}
	// sePlayer is never GCed as long as it plays.
	sePlayer.Play()
}
