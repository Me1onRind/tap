package ui

import (
	"sync"
	"tap/rpc_client"
	"time"
)

type seeker struct {
	mutex   sync.Mutex
	running bool
	tick    bool
	w       *Window

	second    int64
	Shielding bool
}

func NewSeek(w *Window) *seeker {
	s := &seeker{w: w}
	s.reset()
	return s
}

func (s *seeker) Handle(second int64, name string, stopFirst bool) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.tick = true

	if len(name) == 0 || second == 0 {
		return
	}
	if s.running {
		s.second = second
	} else {
		s.Shielding = true
		if stopFirst {
			rpc_client.PlayOrPause(name)
		}
		s.running = true
		go s.pipe(name)
	}
}

func (s *seeker) reset() {
	s.running = false
	s.tick = false
	s.second = 0
	s.Shielding = false
}

func (s *seeker) pipe(name string) {
	for {
		time.Sleep(time.Millisecond * 300)
		s.mutex.Lock()
		if !s.tick {
			// todo send
			rpc_client.SeekAudioFile(s.second)
			s.reset()
			rpc_client.PlayOrPause(name)
			s.mutex.Unlock()
			break
		}

		s.tick = false
		s.mutex.Unlock()
	}
}
