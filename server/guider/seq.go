package guider

import (
	"container/list"
	"log"
	//"math/rand"
	"tap/backend"
	//"time"
)

type seq struct {
	provider backend.Provider

	audioPathList *list.List
	it            *list.Element
}

func (s *seq) NextAudioPath() string {
	next := s.it.Next()
	if next == nil {
		next = s.audioPathList.Front()
	}
	s.it = next
	audioPath := next.Value.(string)
	return audioPath
}

func (s *seq) SetCurrAudioPath(audioPath string) {
	if s.audioPathList.Len() == 0 {
		items, err := s.provider.ListAll()
		if err != nil {
			log.Println(err)
			return
		}
		s.audioPathList = list.New()
		for _, v := range items {
			e := s.audioPathList.PushBack(v)
			if v == audioPath {
				s.it = e
			}
		}
	}
}
