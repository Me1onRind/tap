package server

import (
	"sync"
	"tap/server/guider"
)

type ProviderType int32

const (
	_                 ProviderType = iota
	LocalPrivoderType              = 1
)

const (
	_ uint32 = iota
	SINGLE_MODE
	RANDOM_MODE
	SEQ_MODE
)

var (
	mutex sync.Mutex
	gd    guider.Guider
)

func init() {
	//worker.AddCallback(func(p *player.PlayerWorker) {
	//audioPath := gd.NextAudioPath()
	//p.Play(audioPath)
	//})
}
