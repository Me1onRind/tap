package server

import (
	"tap/player"
)

func PlayOrPause(name string) (*player.AudioInfo, error) {
	mutex.Lock()
	defer mutex.Unlock()

	currInfo, err := worker.CurrAudioInfo()
	if err != nil {
		return nil, err
	}

	audiopath, err := provider.Filepath(name)
	if err != nil {
		return nil, err
	}

	if currInfo.Pathinfo == audiopath && currInfo.Status == player.PLAY {
		worker.Pause()
		currInfo.Status = player.PAUSE
		currName = name
		return currInfo, nil
	} else {
		if err := worker.Play(audiopath); err != nil {
			return nil, err
		}
		currName = name
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

func SetVolume(volume float32) {
	mutex.Lock()
	defer mutex.Unlock()
	worker.SetVolume(volume)
}
