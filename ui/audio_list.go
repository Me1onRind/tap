package ui

import (
	"fmt"
	"github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"tap/server"
)

type audioList struct {
	self   *widgets.List // play status
	window *Window

	rowsChan     chan []string
	playNameChan chan string

	playName string
}

func newAudioList(window *Window) *audioList {
	a := &audioList{
		self:   widgets.NewList(),
		window: window,

		rowsChan:     make(chan []string, _CHANNEL_SIZE),
		playNameChan: make(chan string, _CHANNEL_SIZE),
	}

	audioListWg := a.self
	audioListWg.TextStyle = termui.NewStyle(termui.Color(204))
	audioListWg.SelectedRowStyle = termui.NewStyle(termui.ColorGreen)
	audioListWg.PaddingLeft = 2
	audioListWg.PaddingTop = 1
	audioListWg.WrapText = false
	a.window.setPersentRect(audioListWg, 0.46, 0.13, 0.4, 0.74)

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
	a.self.Rows = a.window.ListAll()
	a.playName = info.Name
	for k, v := range a.self.Rows {
		if v == a.playName {
			a.self.SelectedRow = k
			break
		}
	}
}

func (a *audioList) Print() {
	a.self.Title = "Audio file list"
	for k, v := range a.self.Rows {
		if v == a.playName {
			a.self.Rows[k] = fmt.Sprintf("[%s](fg:yellow)", v)
			termui.Render(a.self)
			a.self.Rows[k] = v
			return
		}
	}
	termui.Render(a.self)

}

func (a *audioList) Cronjob() {
	for {
		select {
		case rows := <-a.rowsChan:
			a.self.Rows = rows
		case playName := <-a.playNameChan:
			a.playName = playName
		}
		a.window.SyncPrint(a.Print)
	}
}

func (a *audioList) HandleEvent(input string) {
	audioListWg := a.self
	switch input {
	case "q", "<C-c>":
		return
	case "j", "<Down>":
		audioListWg.ScrollDown()
	case "k", "<Up>":
		audioListWg.ScrollUp()
	case "<C-d>":
		audioListWg.ScrollHalfPageDown()
	case "<C-u>":
		audioListWg.ScrollHalfPageUp()
	case "<Enter>":
		a.playOrPause()
	case "<Space>":
		a.playOrPause()
	}
}

func (a *audioList) NotifyRowsChange(rows []string) {
	a.rowsChan <- rows
}

func (a *audioList) NotifyPlayNameChange(name string) {
	a.playNameChan <- name
}
func (a *audioList) playOrPause() {
	if a.self.SelectedRow >= len(a.self.Rows) {
		return
	}
	a.window.PlayOrPause(a.self.Rows[a.self.SelectedRow])
}
