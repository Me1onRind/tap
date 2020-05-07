package server

import (
	"fmt"
	"tap/player"
)

func PlayOrPause(audiopath string) (*player.AudioInfo, error) {
	mutex.Lock()
	defer mutex.Unlock()

	currInfo, err := worker.CurrAudioInfo()
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
		fmt.Println(audiopath)
		currInfo, err = worker.CurrAudioInfo()
		if err != nil {
			return nil, err
		}
		return currInfo, nil
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

func Seek(second int64) {
	mutex.Lock()
	defer mutex.Unlock()
	worker.Seek(second)
}
