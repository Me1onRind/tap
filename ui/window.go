package ui

import (
	"github.com/gizak/termui/v3"
	"sync"
	"tap/rpc_client"
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

	playStatus       *playStatus
	audioList        *audioList
	dirList          *dirList
	volumeController *volumeController
	searchInput      *searchInput
	helpBox          *helpBox
	rightHelpBox     *rightHelpBox
	output           *output

	initPrinters []initPrinter
	printers     []printer
	tabItems     []item
	focus        item
	tabIndex     int

	levelOffset float64
	mutex       sync.Mutex
}

func NewWindow() *Window {
	w := &Window{
		levelOffset: 0.00,
	}
	return w
}

func (w *Window) Init() {
	termui.Init()
	w.MaxX, w.MaxY = termui.TerminalDimensions()

	w.initMember()
	w.startPrint()

	// cronjob
	go w.playStatus.Cronjob()
	go w.audioList.Cronjob()

	// tab term list
	w.tabItems = append(w.tabItems, w.audioList)
	w.tabItems = append(w.tabItems, w.dirList)

	w.ChoseItem(w.audioList)

	go rpc_client.Subscribe(func(info *server.PlayAudioInfo) {
		w.playStatus.Notify(info)
		w.audioList.NotifyAudioPathChange(info.Name)
	})

	uiEvents := termui.PollEvents()
	for {
		e := <-uiEvents
		switch e.ID {
		case "<Tab>":
			w.nextItem()
		case "<C-f>":
			w.ChoseItem(w.searchInput)
		case "<C-a>":
			w.ChoseItem(w.audioList)
		case "<Up>":
			w.volumeController.Up()
		case "<Down>":
			w.volumeController.Down()
		case "<C-n>":
			w.playStatus.ChangeLoopMode()
		case "<Left>":
			w.playStatus.SeekAudioFile(-2)
		case "<Right>":
			w.playStatus.SeekAudioFile(2)
		case "<C-c>", "<C-q>", "<Escape>":
			return
		default:
			w.focus.HandleEvent(e.ID)
			w.focus.Print()
		}
	}
}

func (w *Window) SyncPrint(print func()) {
	w.mutex.Lock()
	print()
	w.mutex.Unlock()
}

func (w *Window) initMember() {
	w.playStatus = newPlayStatus(w)
	w.volumeController = newVolumeController(w)
	w.audioList = newAudioList(w)

	w.initPrinters = append(w.initPrinters, w.playStatus)
	w.initPrinters = append(w.initPrinters, w.volumeController)
	w.initPrinters = append(w.initPrinters, w.audioList)

	w.dirList = newDirList(w)
	w.searchInput = newSearchInput(w)
	w.helpBox = newHelpBox(w)
	w.rightHelpBox = newRightHelpBox(w)
	w.output = newOutput(w)

	w.printers = append(w.printers, w.dirList)
	w.printers = append(w.printers, w.searchInput)
	w.printers = append(w.printers, w.helpBox)
	w.printers = append(w.printers, w.rightHelpBox)
	w.printers = append(w.printers, w.output)
}

func (w *Window) startPrint() {
	info := rpc_client.PlayStatus()
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
	w.focus.Leave()
	if w.tabIndex == len(w.tabItems)-1 {
		w.tabIndex = 0
	} else {
		w.tabIndex++
	}
	w.focus = w.tabItems[w.tabIndex]
	w.focus.Entry()
}

func (w *Window) ChoseItem(it item) {
	defer func() {
		if w.focus != nil {
			w.focus.Leave()
		}
		w.focus = it
		w.focus.Entry()
	}()

	for k, v := range w.tabItems {
		if v == it {
			w.tabIndex = k
			return
		}
	}
}

func (w *Window) entry() {
	it := w.tabItems[w.tabIndex]
	w.rightHelpBox.UpdateText(it.WidgetKeys())
	it.Entry()

}

func (w *Window) Close() {
	termui.Close()
}

func (w *Window) GetMax() (float64, float64) {
	return float64(w.MaxX), float64(w.MaxY)
}

func (w *Window) GetOutput() rpc_client.Output {
	return w.output
}
