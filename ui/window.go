package ui

import (
	"github.com/gizak/termui/v3"
	//"log"
	//"math"
	"sync"
	"tap/server"
)

const (
	_CHANNEL_SIZE = 160
)

type item interface {
	Entry()
	Leave()
	HandleEvent(intput string)
	Print()
}

type initPrinter interface {
	InitPrint(info *server.PlayAudioInfo)
}

type Window struct {
	MaxX int
	MaxY int

	playerClient server.PlayClient
	ps           *playStatus
	al           *audioList
	dl           *dirList
	vc           *volumeController
	si           *searchInput
	hb           *helpBox
	rhb          *rightHelpBox
	op           *output

	initPrinters []initPrinter
	tabItems     []item
	tabIndex     int

	levelOffset float64
	mutex       sync.Mutex
}

func NewWindow(rpcClient server.PlayClient) *Window {
	w := &Window{
		playerClient: rpcClient,
		levelOffset:  0.00,
	}
	return w
}

func (w *Window) Init() {
	termui.Init()
	w.MaxX, w.MaxY = termui.TerminalDimensions()

	w.initMember()
	w.startPrint()

	// cronjob
	go w.ps.Cronjob()
	go w.al.Cronjob()

	// tab term list
	w.tabItems = append(w.tabItems, w.al)
	w.tabItems = append(w.tabItems, w.si)
	w.tabItems[0].Entry()

	go w.subscribe()

	uiEvents := termui.PollEvents()
	for {
		e := <-uiEvents
		switch e.ID {
		case "<Tab>":
			w.nextItem()
		case "<C-f>":
			w.choseItem(w.si)
		case "<C-a>":
			w.choseItem(w.al)
		case "<C-k>":
			w.vc.Up()
		case "<C-j>":
			w.vc.Down()
		case "<C-n>":
			w.ps.ChangeLoopMode()
		case "<Left>":
			w.ps.SeekAudioFile(-2)
		case "<Right>":
			w.ps.SeekAudioFile(2)
		case "<C-c>", "<C-q>", "<Escape>":
			return
		default:
			w.tabItems[w.tabIndex].HandleEvent(e.ID)
			w.tabItems[w.tabIndex].Print()
		}
	}
}

func (w *Window) SyncPrint(print func()) {
	w.mutex.Lock()
	print()
	w.mutex.Unlock()
}

func (w *Window) initMember() {
	w.ps = newPlayStatus(w)
	w.al = newAudioList(w)
	w.dl = newDirList(w)
	w.vc = newVolumeController(w)
	w.si = newSearchInput(w)
	w.hb = newHelpBox(w)
	w.rhb = newRightHelpBox(w)
	w.op = newOutput(w)

	w.initPrinters = append(w.initPrinters, w.ps)
	w.initPrinters = append(w.initPrinters, w.vc)
	w.initPrinters = append(w.initPrinters, w.al)

}

func (w *Window) startPrint() {
	info := w.PlayStatus()
	if info == nil {
		return
	}

	for _, v := range w.initPrinters {
		v.InitPrint(info)
	}

	w.si.Print()
	w.hb.Print()
	w.rhb.Print()
	w.op.Print()
	w.dl.Print()
}

func (w *Window) nextItem() {
	w.tabItems[w.tabIndex].Leave()
	if w.tabIndex == len(w.tabItems)-1 {
		w.tabIndex = 0
	} else {
		w.tabIndex++
	}
	w.tabItems[w.tabIndex].Entry()
}

func (w *Window) choseItem(it item) {
	for k, v := range w.tabItems {
		if v == it {
			w.tabItems[w.tabIndex].Leave()
			w.tabIndex = k
			it.Entry()
			return
		}
	}
}

func (w *Window) Close() {
	termui.Close()
}

func (w *Window) GetMax() (float64, float64) {
	return float64(w.MaxX), float64(w.MaxY)
}
