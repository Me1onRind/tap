package server

import (
	"container/list"
	"log"
	"math/rand"
	"sync"
	"tap/backend"
	"tap/player"
	"time"
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
	worker       *player.PlayerWorker
	provider     backend.Provider
	providerType ProviderType
	providerName string
	mutex        sync.Mutex

	// status
	mode     uint32
	sf       *shufflePathList
	currName string
)

func init() {
	worker = player.NewPlayerWorker()
	worker.AddCallback(func(p *player.PlayerWorker) {
		var next string
		if mode == SINGLE_MODE {
			next = currName
		} else if mode == RANDOM_MODE {
			next = sf.next()
		} else if mode == SEQ_MODE {
			next = seqNext(p.CurrAudiopath)
		}
		info, _ := PlayOrPause(next)
		if p != nil {
			ps.push(info)
		}
	})

	mode = RANDOM_MODE
	sf = &shufflePathList{}
}

type shufflePathList struct {
	pathList  *list.List
	currAudio *list.Element
	begin     backend.AudioItem
}

func (s *shufflePathList) buildShufflePathList(currAudio string) {
	pathList := list.New()
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	items, err := provider.ListAll()
	if err != nil {
		log.Println(err)
		return
	}
	for i := 0; i < len(items); i++ {
		index := r.Intn(len(items)-i) + i
		items[i], items[index] = items[index], items[i]
	}

	for _, v := range items {
		if currAudio == v.Name() {
			s.currAudio = pathList.PushBack(v)
			s.begin = v
		} else {
			pathList.PushBack(v)
		}
	}

	s.pathList = pathList
}

func (s *shufflePathList) next() string {
	if s.pathList == nil {
		s.buildShufflePathList(currName)
	}
	log.Printf("currAudio: %s\n", s.currAudio.Value.(backend.AudioItem).Name())

	next := s.currAudio.Next()
	if next == nil {
		next = s.pathList.Front()
	}

	nextItem := next.Value.(backend.AudioItem)
	log.Printf("nextItem: %s\n", nextItem.Name())
	if nextItem == s.begin {
		s.buildShufflePathList(nextItem.Name())
	}
	s.currAudio = next

	return nextItem.Name()
}

func seqNext(currAudiopath string) string {
	items, err := provider.ListAll()
	if err != nil {
		log.Println(err)
		return ""
	}
	if len(items) == 0 {
		return ""
	}

	index := 0
	for index < len(items)-1 {
		if items[index].Name() == currName {
			return items[index+1].Name()
		}
		index++
	}
	return items[0].Name()
}
