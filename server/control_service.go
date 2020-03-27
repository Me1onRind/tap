package server

import (
	//"sync"
	//"tap/backend"
	"tap/player"
)

func Status() (*player.AudioInfo, error) {
	mutex.Lock()
	defer mutex.Unlock()

	return worker.CurrAudioInfo()
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

func Search(input string) ([]string, error) {
	list, err := provider.Search(input)
	if err != nil {
		return nil, err
	}
	var ret []string
	for _, v := range list {
		ret = append(ret, v.Name())
	}
	return ret, nil
}

func SetPlayMoel(m uint32) {
	mode = m
}
