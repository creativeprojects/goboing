package main

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/vorbis"
	"github.com/markbates/pkger"
)

// AudioPlayer represents the current audio state.
type AudioPlayer struct {
	audioContext *audio.Context
	audioPlayer  *audio.Player
	volume128    int
}

func NewAudioPlayer(audioContext *audio.Context) (*AudioPlayer, error) {
	type audioStream interface {
		io.ReadSeeker
		Length() int64
	}

	const bytesPerSample = 4 // TODO: This should be defined in audio package

	var s audioStream
	var err error
	file, err := pkger.Open("/music/theme.ogg")
	if err != nil {
		return nil, err
	}
	defer file.Close()
	theme, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	s, err = vorbis.Decode(audioContext, bytes.NewReader(theme))
	if err != nil {
		return nil, err
	}

	p, err := audio.NewPlayer(audioContext, audio.NewInfiniteLoop(s, s.Length()))
	if err != nil {
		return nil, err
	}
	player := &AudioPlayer{
		audioContext: audioContext,
		audioPlayer:  p,
		volume128:    32,
	}
	player.audioPlayer.SetVolume(float64(player.volume128) / 128)
	player.audioPlayer.Play()

	return player, nil
}

func (p *AudioPlayer) Close() error {
	return p.audioPlayer.Close()
}

// PlaySE plays a sound effect.
func PlaySE(audioContext *audio.Context, bs []byte) {
	if bs == nil || len(bs) == 0 {
		log.Printf("cannot play empty sound")
		return
	}
	sePlayer := audio.NewPlayerFromBytes(audioContext, bs)
	// sePlayer is never GCed as long as it plays.
	sePlayer.Play()
}
