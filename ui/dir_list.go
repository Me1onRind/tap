package ui

import (
	"github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	//"tap/server"
)

type dirList struct {
	self   *widgets.List // play status
	window *Window

	rowsChan     chan []string
	playNameChan chan string

	playName string
}

func newDirList(window *Window) *dirList {
	d := &dirList{
		self:   widgets.NewList(),
		window: window,

		rowsChan:     make(chan []string, _CHANNEL_SIZE),
		playNameChan: make(chan string, _CHANNEL_SIZE),
	}

	maxX, _ := d.window.GetMax()

	d.self.SetRect(_HELP_BOX_WIDTH,
		d.window.MaxY-_COUNT_DOWN_HEIGHT-_PLAY_STATUS_HEIGHT-_HELP_BOX_HEIGHT,
		int(maxX*(_VOLUME_WIDTH+_PLAY_STATUS_WIDTH)),
		d.window.MaxY-_COUNT_DOWN_HEIGHT-_PLAY_STATUS_HEIGHT)
	d.self.Title = "Directory list"
	//a.window.setPersentRect(audioListWg, 0.46, 0.13, 0.4, 0.74)

	return d
}

func (d *dirList) Entry() {
	d.self.BorderStyle.Fg = termui.ColorGreen
	d.Print()
}

func (d *dirList) Leave() {
	d.self.BorderStyle.Fg = termui.ColorWhite
	d.Print()
}

func (d *dirList) Print() {
	termui.Render(d.self)
}

func (d *dirList) HandleEvent(input string) {
	switch input {
	}
}
