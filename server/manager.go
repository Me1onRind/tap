package server

import (
	//"sync"
	"tap/backend"
	"tap/player"
	"tap/server/guider"
)

var m *manager

func init() {
	playWorker := player.NewPlayerWorker()
	m = &manager{
		PlayWorker: playWorker,
	}

	playWorker.AddCallback(func(p *player.PlayerWorker) {
		next := m.nextAudioGuider.NextAudioPath()
		info, _ := m.PlayOrPause(next)
		if info != nil {
			ps.push(info)
		}
	})

}

type manager struct {
	ProviderType int32
	Provider     backend.Provider

	CurrDir string
	Dirs    []string

	PlayMode   uint32
	PlayWorker *player.PlayerWorker

	nextAudioGuider guider.Guider
}

func (m *manager) Init(provider backend.Provider, providerType int32, dirs []string, mode uint32) {
	m.Provider = provider
	m.ProviderType = providerType
	m.Dirs = dirs
	m.CurrDir = dirs[0]
	m.SetPlayMode(mode)
}

func (m *manager) SetPlayMode(mode uint32) {
	if mode == m.PlayMode {
		return
	}
	currInfo, err := m.PlayWorker.CurrAudioInfo()
	if err != nil {
		return
	}
	m.nextAudioGuider = guider.NewGuider(guider.Mode(mode), m.Provider, &m.CurrDir)
	m.nextAudioGuider.SetCurrAudioPath(currInfo.Pathinfo)
	m.PlayMode = mode

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

	if len(currInfo.Pathinfo) == 0 {
		m.nextAudioGuider.SetCurrAudioPath(audioPath)
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
