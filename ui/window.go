package ui

import (
	"github.com/gizak/termui/v3"
	//"log"
	"math"
	"sync"
	"tap/server"
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
	MaxX float64
	MaxY float64

	playerClient server.PlayClient
	ps           *playStatus
	al           *audioList
	vc           *volumeController
	si           *searchInput

	tabItems []item
	tabIndex int

	initPrinters []initPrinter

	levelOffset float64

	mutex sync.Mutex
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
	maxX, maxY := termui.TerminalDimensions()
	w.MaxX = float64(maxX)
	w.MaxY = float64(maxY)

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
		case "<C-n>":
			w.ps.ChangeLoopMode()
		case "<C-j>":
			w.vc.Down()
		case "<C-k>":
			w.vc.Up()
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
	w.vc = newVolumeController(w)
	w.si = newSearchInput(w)

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

func (w *Window) Close() {
	termui.Close()
}

func (w *Window) setPersentRect(block termui.Drawable, offsetX, offsetY, width, height float64) {
	x0 := int(math.Ceil(w.MaxX*offsetX + w.levelOffset*w.MaxX))
	y0 := int(w.MaxY * offsetY)
	x1 := x0 + int(math.Ceil(w.MaxX*width))
	y1 := y0 + int(w.MaxY*height)
	block.SetRect(x0, y0, x1, y1)
}

func (w *Window) setPersentRectWithFixedHeight(block termui.Drawable, offsetX, offsetY, width float64, height int) {
	x0 := int(math.Ceil(w.MaxX*offsetX + w.levelOffset*w.MaxX))
	y0 := int(w.MaxY * offsetY)
	x1 := x0 + int(math.Ceil(w.MaxX*width))
	y1 := y0 + height
	block.SetRect(x0, y0, x1, y1)
}
