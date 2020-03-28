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
		afterPlaySucc(currInfo)
		return currInfo, nil
	} else {
		if err := worker.Play(audiopath); err != nil {
			return nil, err
		}
		currName = name
		currInfo, err = worker.CurrAudioInfo()
		if err != nil {
			return nil, err
		}
		afterPlaySucc(currInfo)
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

func Forward(second uint32) {
	info, _ := worker.Forward(second)
	ps.push(info)
}

func Rewind(second uint32) {
	info, _ := worker.Rewind(second)
	ps.push(info)
}

func afterPlaySucc(info *player.AudioInfo) {
	ps.push(info)

	if mode == RANDOM_MODE && info.Status == player.PLAY {
		sf.pathList = nil
	}
}
