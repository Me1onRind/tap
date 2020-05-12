package ui

import (
	"fmt"
	"github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"path/filepath"
	"tap/rpc_client"
	"tap/server"
)

type audioList struct {
	self   *widgets.List // play status
	window *Window
	rows   []string

	rowsChan      chan []string
	audioPathChan chan string

	audioPath string
}

func newAudioList(window *Window) *audioList {
	a := &audioList{
		self:   widgets.NewList(),
		window: window,

		rowsChan:      make(chan []string, _CHANNEL_SIZE),
		audioPathChan: make(chan string, _CHANNEL_SIZE),
	}

	audioListWg := a.self
	audioListWg.TextStyle = termui.NewStyle(termui.Color(204))
	audioListWg.SelectedRowStyle = termui.NewStyle(termui.ColorGreen)
	audioListWg.PaddingLeft = 2
	audioListWg.PaddingTop = 1
	audioListWg.WrapText = false
	a.self.Title = "Audio file list"
	maxX, _ := a.window.GetMax()
	audioListWg.SetRect(int(maxX*(_VOLUME_WIDTH+_PLAY_STATUS_WIDTH)), _SEARCH_HEIGHT,
		int(maxX*(_VOLUME_WIDTH+_PLAY_STATUS_WIDTH+_LIST_WIDTH)), a.window.MaxY-_COUNT_DOWN_HEIGHT)
	//a.window.setPersentRect(audioListWg, 0.46, 0.13, 0.4, 0.74)

	return a
}

func (a *audioList) Entry() {
	a.self.BorderStyle.Fg = termui.ColorGreen
	a.Print()
}

func (a *audioList) Leave() {
	a.self.BorderStyle.Fg = termui.ColorWhite
	a.Print()
}

func (a *audioList) InitPrint(info *server.PlayAudioInfo) {
	a.rows = rpc_client.ListAll()
	a.audioPath = info.Name
	for k, v := range a.self.Rows {
		if v == a.audioPath {
			a.self.SelectedRow = k
			break
		}
	}
	a.Print()
}

func (a *audioList) HandleEvent(input string) {
	audioListWg := a.self
	switch input {
	case "j", "<Down>":
		audioListWg.ScrollDown()
	case "k", "<Up>":
		audioListWg.ScrollUp()
	case "<C-d>":
		audioListWg.ScrollHalfPageDown()
	case "<C-u>":
		audioListWg.ScrollHalfPageUp()
	case "<Enter>", "<Space>":
		a.playOrPause()
	}
}

func (a *audioList) WidgetKeys() string {
	return "j <Down> Select next\n" +
		"k <Up>   Select prev\n" +
		"<C-d>    Page Down\n" +
		"<C-u>    Page Up\n" +
		"<Enter>  Play or Pause\n" +
		"<Space>  Play or Pause\n"
}

func (a *audioList) Print() {
	var rows []string
	for _, v := range a.rows {
		if v == a.audioPath {
			rows = append(rows, fmt.Sprintf("[%s](fg:yellow)", filepath.Base(v)))
		} else {
			rows = append(rows, filepath.Base(v))
		}
	}
	a.self.Rows = rows
	termui.Render(a.self)

}

func (a *audioList) Cronjob() {
	for {
		select {
		case rows := <-a.rowsChan:
			a.self.Rows = rows
		case audioPath := <-a.audioPathChan:
			a.audioPath = audioPath
		}
		a.window.SyncPrint(a.Print)
	}
}

func (a *audioList) NotifyRowsChange(rows []string) {
	a.rowsChan <- rows
}

func (a *audioList) NotifyAudioPathChange(audioPath string) {
	a.audioPathChan <- audioPath
}
func (a *audioList) playOrPause() {
	if a.self.SelectedRow >= len(a.self.Rows) {
		return
	}
	rpc_client.PlayOrPause(a.rows[a.self.SelectedRow])
}
