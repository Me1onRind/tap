package server

import (
	//"sync"
	"tap/backend"
	"tap/player"
)

var m *manager

func init() {
	m = &manager{
		PlayWorker: player.NewPlayerWorker(),
	}
}

type manager struct {
	ProviderType int32
	Provider     backend.Provider

	PlayMode uint32
	CurrDir  string
	Dirs     []string

	PlayWorker *player.PlayerWorker
}

func (m *manager) Init(provider backend.Provider, providerType int32, dirs []string) {
	m.Provider = provider
	m.ProviderType = providerType
	m.Dirs = dirs
	m.CurrDir = dirs[0]
}

func (m *manager) ListAll() ([]string, error) {
	return m.Provider.ListAll(m.CurrDir)
}

func (m *manager) Search(input string) ([]string, error) {
	return m.Provider.Search(input, m.CurrDir)
}

func (m *manager) PlayOrPause(audioPath string) (*player.AudioInfo, error) {
	currInfo, err := m.PlayWorker.CurrAudioInfo()
	if err != nil {
		return nil, err
	}

	if currInfo.Pathinfo == audioPath && currInfo.Status == player.PLAY {
		m.PlayWorker.Pause()
		currInfo.Status = player.PAUSE
		return currInfo, nil
	} else {
		if err := m.PlayWorker.Play(audioPath); err != nil {
			return nil, err
		}
		currInfo, err = m.PlayWorker.CurrAudioInfo()
		if err != nil {
			return nil, err
		}
		return currInfo, nil
	}
}
