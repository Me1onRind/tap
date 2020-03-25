package ui

import (
	"context"
	"github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"tap/server"
)

type audioList struct {
	self   *widgets.List // play status
	window *Window
}

func newAudioList(window *Window) *audioList {
	a := &audioList{
		self:   widgets.NewList(),
		window: window,
	}

	audioListWg := a.self
	audioListWg.TextStyle = termui.NewStyle(termui.ColorYellow)
	audioListWg.PaddingLeft = 2
	audioListWg.PaddingTop = 1
	audioListWg.WrapText = false

	a.window.setPersentRect(audioListWg, 0.46, 0.05, 0.4, 0.82)
	return a
}

func (a *audioList) print() {
	audioListWg := a.self

	audioListWg.Title = "Audio file list"
	res, _ := a.window.playerClient.ListAll(context.Background(), &server.Empty{})

	audioListWg.Rows = res.GetNames()
	a.window.syncPrint(audioListWg)
	uiEvents := termui.PollEvents()
	for {
		e := <-uiEvents
		a.handleEvent(&e)
		termui.Render(audioListWg)
		if e.ID == "q" {
			break
		}
	}
}

func (a *audioList) handleEvent(e *termui.Event) {
	audioListWg := a.self
	switch e.ID {
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
		audioListWg.ScrollPageDown()
	case "<C-b>":
		audioListWg.ScrollPageUp()
	case "<Enter>":
		a.PlayOrPause()
	case "<Space>":
		a.PlayOrPause()
	default:
	}
}

func (a *audioList) PlayOrPause() {
	res := a.window.PlayOrPause(a.self.SelectedRow)
	if res != nil {
		a.window.ps.flushForce <- res
	}
}
