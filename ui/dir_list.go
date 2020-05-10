package ui

import (
	"fmt"
	"github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"path/filepath"
	//"tap/server"
)

type dirList struct {
	self   *widgets.List // play status
	window *Window

	dirs    []string
	currDir string
	//rowsChan     chan []string
	//playNameChan chan string

	//playName string
}

func newDirList(w *Window) *dirList {
	d := &dirList{
		self:   widgets.NewList(),
		window: w,
	}

	d.self.TextStyle = termui.NewStyle(termui.Color(204))
	d.self.SelectedRowStyle = termui.NewStyle(termui.ColorGreen)

	maxX, _ := d.window.GetMax()

	d.self.SetRect(_HELP_BOX_WIDTH,
		d.window.MaxY-_COUNT_DOWN_HEIGHT-_PLAY_STATUS_HEIGHT-_HELP_BOX_HEIGHT,
		int(maxX*(_VOLUME_WIDTH+_PLAY_STATUS_WIDTH)),
		d.window.MaxY-_COUNT_DOWN_HEIGHT-_PLAY_STATUS_HEIGHT)
	d.self.Title = "Directory list"
	providerInfo := d.window.Provider()
	if providerInfo != nil {
		d.dirs = providerInfo.Dirs
		d.currDir = providerInfo.CurrDir
	}
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
	var rows []string
	for _, v := range d.dirs {
		if v == d.currDir {
			rows = append(rows, fmt.Sprintf("[%s](fg:yellow)", filepath.Base(v)))
		} else {
			rows = append(rows, filepath.Base(v))
		}
	}
	d.self.Rows = rows
	termui.Render(d.self)
}

func (d *dirList) WidgetKeys() string {
	return "j <Down> Select next\n" +
		"k <Up>   Select prev\n" +
		"<Entry>  Chose Dir\n" +
		"<Space>  Chose Dir\n"
}

func (d *dirList) HandleEvent(input string) {
	switch input {
	case "j", "<Down>":
		d.self.ScrollDown()
	case "k", "<Up>":
		d.self.ScrollUp()
	case "<Enter>", "<Space>":
		d.currDir = d.dirs[d.self.SelectedRow]
		d.window.SetDir(d.currDir)
		rows := d.window.ListAll()
		d.window.audioList.NotifyRowsChange(rows)
		d.Print()
	}
}
