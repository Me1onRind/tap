package ui

import (
	//"context"
	"fmt"
	"github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"path/filepath"
	"tap/player"
	"tap/server"
	"time"
)

type playStatus struct {
	self      *widgets.Paragraph // play status
	window    *Window
	progress  *widgets.Gauge
	countDown *widgets.Paragraph

	flushForce chan *server.PlayAudioInfo

	AudioName   string
	Duration    uint32
	CurrPro     uint32
	SampleRate  uint32
	Status      uint32
	StatusLabel string
}

func newPlayStatus(w *Window) *playStatus {
	p := &playStatus{
		self:        widgets.NewParagraph(),
		progress:    widgets.NewGauge(),
		countDown:   widgets.NewParagraph(),
		window:      w,
		flushForce:  make(chan *server.PlayAudioInfo, 10),
		Status:      0,
		StatusLabel: "Stop",
	}

	p.self.Title = "Status"
	p.self.PaddingLeft = 3

	p.progress.Label = " "
	p.countDown.Text = " 00:00/00:00"

	p.window.setPersentRect(p.self, 0.14, 0.61, 0.315, 0.27)
	p.window.setPersentRect(p.progress, 0.07, 0.87, 0.685, 0.09)
	p.window.setPersentRect(p.countDown, 0.758, 0.87, 0.10, 0.09)
	return p
}

func (p *playStatus) printStatus() {
	pg := p.self
	pg.Text = p.text()
	p.window.syncPrint(pg)
}

func (p *playStatus) printPro() {
	p.window.syncPrint(p.countDown)
	p.window.syncPrint(p.progress)
}

func (p *playStatus) asyncPrint() {
	p.printPro()
	p.printStatus()
	ticker := time.NewTicker(time.Second)
	for {
		select {
		case <-ticker.C:
			if p.Status == player.PLAY && p.CurrPro < p.Duration {
				p.CurrPro++
			}
			p.updateProgress()
			p.printPro()
		case a := <-p.flushForce:
			p.Status = a.GetStatus()
			if a.GetPathinfo() != "" {
				p.AudioName = filepath.Base(a.GetPathinfo())
			}
			p.Duration = a.GetDuration()
			p.CurrPro = a.GetCurr()
			p.SampleRate = a.GetSampleRate()

			p.updateProgress()
			ticker.Stop()
			ticker = time.NewTicker(time.Second)
			p.printPro()
			p.printStatus()
		}
	}
}

func (p *playStatus) updateProgress() {
	if p.Duration == 0 {
		p.progress.Percent = 0
	} else {
		p.progress.Percent = int(100 * p.CurrPro / p.Duration)
	}
	p.countDown.Text = fmt.Sprintf(" %s/%s", formatDuration(p.CurrPro), formatDuration(p.Duration))
	switch p.Status {
	case player.PLAY:
		p.StatusLabel = "Playing ▶️  "
		p.self.BorderStyle = termui.NewStyle(termui.ColorYellow)
	case player.PAUSE:
		p.StatusLabel = "Pause ⏸  "
		p.self.BorderStyle = termui.NewStyle(termui.ColorWhite)
	default:
		p.StatusLabel = "Stop ⏹  "
		p.self.BorderStyle = termui.NewStyle(termui.ColorWhite)
	}
}

func (p *playStatus) text() string {
	return fmt.Sprintf(
		"\n%s\n"+
			"\nAudio:      %s\n"+
			"\nDuration:   %s\n"+
			"\nSampleRate: %.1f kHz\n",
		p.StatusLabel, p.AudioName, formatDuration(p.Duration), float64(p.SampleRate)/1000.0)
}

func formatDuration(t uint32) string {
	return fmt.Sprintf("%02d:%02d", t/60, t%60)
}
