package server

import (
	"sync"
	"tap/backend"
	"tap/player"
)

type ProviderType int32

const (
	_                 ProviderType = iota
	LocalPrivoderType              = 1
)

var (
	worker       *player.PlayerWorker
	provider     backend.Provider
	providerType ProviderType
	providerName string
	mutex        sync.Mutex
)

func init() {
	worker = player.NewPlayerWorker()
}

func Status() (*player.AudioInfo, error) {
	mutex.Lock()
	defer mutex.Unlock()

	return worker.CurrAudioInfo()
}

func PlayOrPause(index int) (*player.AudioInfo, error) {
	mutex.Lock()
	defer mutex.Unlock()

	currInfo, err := worker.CurrAudioInfo()
	if err != nil {
		return nil, err
	}

	audiopath, err := provider.Filepath(index)
	if err != nil {
		return nil, err
	}

	if currInfo.Pathinfo == audiopath && currInfo.Status == player.PLAY {
		worker.Pause()
		currInfo.Status = player.PAUSE
		return currInfo, nil
	} else {
		if err := worker.Play(audiopath); err != nil {
			return nil, err
		}
		return worker.CurrAudioInfo()
	}

}

func Pause() {
	mutex.Lock()
	defer mutex.Unlock()
	worker.Pause()
}

func Stop() {
	mutex.Lock()
	defer mutex.Unlock()
	worker.Stop()
}

func ListAll() ([]string, error) {
	list, err := provider.ListAll()
	if err != nil {
		return nil, err
	}
	var ret []string
	for _, v := range list {
		ret = append(ret, v.Name())
	}
	return ret, nil
}
