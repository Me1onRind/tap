package guider

import (
	"tap/backend"
)

type Guider interface {
	NextAudioPath() string
	PreAudioPath() string
	SetCurrAudioPath(audioPath string)
}

type Mode uint8

const (
	_ Mode = iota
	RANDOM
	SEQ
	SINGLE
)

func NewGuider(mode Mode, provider backend.Provider, currDir *string) Guider {
	if mode == RANDOM {
		return newRandom(provider, currDir)
	}

	if mode == SEQ {
		return newSeq(provider, currDir)
	}

	if mode == SINGLE {
		return NewSingle()
	}

	return nil
}

type Single struct {
	audioPath string
}

func NewSingle() *Single {
	return &Single{}
}

func (s *Single) NextAudioPath() string {
	return s.audioPath
}

func (s *Single) PreAudioPath() string {
	return ""
}

func (s *Single) SetCurrAudioPath(audioPath string) {
	s.audioPath = audioPath
}
