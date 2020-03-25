package player

import (
	"testing"
	"time"
)

func Test_Play_One(t *testing.T) {
	audiopath := "/Users/me1onrind/Downloads/partyboobytrap-Like Ships.mp3"
	p := NewPlayerWorker()
	err := p.Play(audiopath)
	t.Log(err)
	time.Sleep(time.Second * 10)
	audioInfo, err := p.CurrAudioInfo()
	t.Log(audioInfo, err)
	select {}
}

func Test_Play_Two(t *testing.T) {
	t.Skip()
	audiopath := "/Users/me1onrind/Downloads/partyboobytrap-Like Ships.mp3"
	p := NewPlayerWorker()
	err := p.Play(audiopath)
	t.Log(err)
	time.Sleep(time.Second * 10)
	audiopath = "/Users/me1onrind/Downloads/syst.mp3"
	err = p.Play(audiopath)
	t.Log(err)
	select {}
}

func Test_Play_Three(t *testing.T) {
}
