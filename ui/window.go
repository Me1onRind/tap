package ui

import (
	"github.com/gizak/termui/v3"
	//"github.com/gizak/termui/v3/widgets"
	sll "github.com/emirpasic/gods/lists/singlylinkedlist"
	//"github.com/emirpasic/gods/utils"
	"math"
	"sync"
	"tap/server"
)

type item interface {
	entry()
	leave()
	handleEvent(intput string)
	print()
}

type Window struct {
	MaxX float64
	MaxY float64

	playerClient server.PlayServerClient
	ps           *playStatus
	al           *audioList
	vc           *volumeController
	si           *searchInput
	mutex        sync.Mutex

	items       *sll.List
	it          sll.Iterator
	currItem    item
	levelOffset float64
}

func NewWindow(rpcClient server.PlayServerClient) *Window {
	w := &Window{
		playerClient: rpcClient,
		levelOffset:  0.00,

		items: sll.New(),
	}
	return w
}

func (w *Window) Init() {
	termui.Init()
	maxX, maxY := termui.TerminalDimensions()
	w.MaxX = float64(maxX)
	w.MaxY = float64(maxY)
	a := w.chceckPlayStatus()
	if a == nil {
		return
	}

	w.ps = newPlayStatus(w)
	w.al = newAudioList(w)
	w.vc = newVolumeController(w, a.GetVolume())
	w.si = newSearchInput(w)

	w.items.Add(w.al)
	w.items.Add(w.si)

	go w.ps.asyncPrint()
	w.ps.flushForce <- a

	w.al.trySelectInit(a.GetName())
	go w.al.asyncPrint()

	w.vc.print()
	w.si.print()
	w.al.print()

	uiEvents := termui.PollEvents()

	v, _ := w.items.Get(0)
	w.currItem = v.(item)
	w.currItem.entry()
	w.it = w.items.Iterator()
	w.it.Next()

	for {
		e := <-uiEvents
		switch e.ID {
		case "<Tab>":
			w.nextItem()
		case "<C-c>", "<C-q>", "<Escape>":
			return
		default:
			w.currItem.handleEvent(e.ID)
			w.currItem.print()
		}
	}
}

func (w *Window) nextItem() {
	w.currItem.leave()
	if !w.it.Next() {
		w.it = w.items.Iterator()
		w.it.Next()
	}
	w.currItem = w.it.Value().(item)
	w.currItem.entry()
}

func (w *Window) syncPrint(d termui.Drawable) {
	w.mutex.Lock()
	termui.Render(d)
	w.mutex.Unlock()
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
