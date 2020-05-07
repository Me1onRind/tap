package guider

import (
	"container/list"
	"log"
	"math/rand"
	"tap/backend"
	"time"
)

type random struct {
	provider      backend.Provider
	audioPathList *list.List
	it            *list.Element

	beginAuthPath string
}

func (r *random) NextAudioPath() string {
	audioPath := r.getNextAudioPath()
	if audioPath == r.beginAuthPath {
		r.SetCurrAudioPath(audioPath)
	}
	return audioPath
}

func (r *random) SetCurrAudioPath(audioPath string) {
	r.beginAuthPath = audioPath
	r.buildAudioPathList()
}

func (r *random) getNextAudioPath() string {
	next := r.it.Next()
	if next == nil {
		next = r.audioPathList.Front()
	}
	r.it = next
	audioPath := next.Value.(string)
	return audioPath
}

func (r *random) buildAudioPathList() {
	r.audioPathList = list.New()
	items, err := r.provider.ListAll()
	if err != nil {
		log.Println(err)
		return
	}

	rd := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < len(items); i++ {
		index := rd.Intn(len(items)-i) + i
		items[i], items[index] = items[index], items[i]
		e := r.audioPathList.PushBack(items[i])
		if items[i] == r.beginAuthPath {
			r.it = e
		}
	}
}
