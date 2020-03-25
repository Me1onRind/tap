package ui

import (
	"github.com/gizak/termui/v3"
	//"github.com/gizak/termui/v3/widgets"
	"math"
	"sync"
	"tap/server"
)

type Window struct {
	MaxX float64
	MaxY float64

	playerClient server.PlayServerClient
	ps           *playStatus
	al           *audioList
	mutex        sync.Mutex

	levelOffset float64
}

func NewWindow(rpcClient server.PlayServerClient) *Window {
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

	w.ps = newPlayStatus(w)
	w.al = newAudioList(w)

	go w.ps.flushPrint()
	a := w.ChceckPlayStatus()
	if a != nil {
		w.ps.flushForce <- a
	}

	w.al.print()

	for e := range termui.PollEvents() {
		if e.Type == termui.KeyboardEvent {
			break
		}
	}
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
