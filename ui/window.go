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
	WidgetKeys() string
}

type initPrinter interface {
	InitPrint(info *server.PlayAudioInfo)
}

type printer interface {
	Print()
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
	printers     []printer
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
	w.tabItems = append(w.tabItems, w.si)
	w.tabItems = append(w.tabItems, w.al)
	w.tabItems = append(w.tabItems, w.dl)
	w.ChoseItem(w.al)

	go w.subscribe()

	uiEvents := termui.PollEvents()
	for {
		e := <-uiEvents
		switch e.ID {
		case "<Tab>":
			w.nextItem()
		case "<C-f>":
			w.ChoseItem(w.si)
		case "<C-a>":
			w.ChoseItem(w.al)
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
	w.vc = newVolumeController(w)
	w.al = newAudioList(w)

	w.initPrinters = append(w.initPrinters, w.ps)
	w.initPrinters = append(w.initPrinters, w.vc)
	w.initPrinters = append(w.initPrinters, w.al)

	w.dl = newDirList(w)
	w.si = newSearchInput(w)
	w.hb = newHelpBox(w)
	w.rhb = newRightHelpBox(w)
	w.op = newOutput(w)

	w.printers = append(w.printers, w.dl)
	w.printers = append(w.printers, w.si)
	w.printers = append(w.printers, w.hb)
	w.printers = append(w.printers, w.rhb)
	w.printers = append(w.printers, w.op)
}

func (w *Window) startPrint() {
	info := w.PlayStatus()
	if info == nil {
		return
	}

	for _, v := range w.initPrinters {
		v.InitPrint(info)
	}

	for _, v := range w.printers {
		v.Print()
	}
}

func (w *Window) nextItem() {
	w.tabItems[w.tabIndex].Leave()
	if w.tabIndex == len(w.tabItems)-1 {
		w.tabIndex = 0
	} else {
		w.tabIndex++
	}
	w.entry()
}

func (w *Window) ChoseItem(it item) {
	for k, v := range w.tabItems {
		if v == it {
			w.tabItems[w.tabIndex].Leave()
			w.tabIndex = k
			w.entry()
			return
		}
	}
}

func (w *Window) entry() {
	it := w.tabItems[w.tabIndex]
	w.rhb.UpdateText(it.WidgetKeys())
	it.Entry()

}

func (w *Window) Close() {
	termui.Close()
}

func (w *Window) GetMax() (float64, float64) {
	return float64(w.MaxX), float64(w.MaxY)
}
