package ui

import (
	"github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

type audioList struct {
	self   *widgets.List // play status
	window *Window

	rowsChan chan []string
}

func newAudioList(window *Window) *audioList {
	a := &audioList{
		self:   widgets.NewList(),
		window: window,

		rowsChan: make(chan []string, 10),
	}

	audioListWg := a.self
	audioListWg.TextStyle = termui.NewStyle(termui.Color(204))
	audioListWg.SelectedRowStyle = termui.NewStyle(termui.ColorGreen)
	audioListWg.PaddingLeft = 2
	audioListWg.PaddingTop = 1
	audioListWg.WrapText = false
	a.window.setPersentRect(audioListWg, 0.46, 0.13, 0.4, 0.74)
	audioListWg.Rows = a.window.listAll()

	return a
}

func (a *audioList) entry() {
	a.self.BorderStyle.Fg = termui.ColorGreen
	a.print()
}

func (a *audioList) leave() {
	a.self.BorderStyle.Fg = termui.ColorWhite
	a.print()
}

func (a *audioList) print() {
	audioListWg := a.self

	audioListWg.Title = "Audio file list"

	a.window.syncPrint(audioListWg)
}

func (a *audioList) handleEvent(input string) {
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
	case "<C-f>":
	case "<Tab>":
	case "<Enter>":
		a.playOrPause()
	case "<Space>":
		a.playOrPause()
	case "<C-j>":
		a.window.vc.down()
	case "<C-k>":
		a.window.vc.up()
	default:
	}
}

func (a *audioList) playOrPause() {
	res := a.window.playOrPause(a.self.Rows[a.self.SelectedRow])
	if res != nil {
		a.window.ps.flushForce <- res
	}
}

func (a *audioList) asyncPrint() {
	for {
		select {
		case rows := <-a.rowsChan:
			a.self.Rows = rows
			a.print()
		}
	}
}

func (a *audioList) trySelectInit(name string) {
	for k, v := range a.self.Rows {
		if v == name {
			a.self.SelectedRow = k
			a.print()
			return
		}
	}
}
