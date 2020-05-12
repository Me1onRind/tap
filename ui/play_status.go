package ui

import (
	"fmt"
	"github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"path/filepath"
	"tap/player"
	"tap/rpc_client"
	"tap/server"
	"tap/server/guider"
	"time"
)

type playStatus struct {
	self      *widgets.Paragraph // play status
	window    *Window
	progress  *widgets.Gauge
	countDown *widgets.Paragraph

	infoChan chan *server.PlayAudioInfo

	skr *seeker

	AudioPath   string
	Status      uint32
	StatusLabel string
	LoopMode    guider.Mode

	Duration int64
	Endline  int64
	CurrPro  int64
}

func newPlayStatus(w *Window) *playStatus {
	p := &playStatus{
		self:        widgets.NewParagraph(),
		progress:    widgets.NewGauge(),
		countDown:   widgets.NewParagraph(),
		window:      w,
		infoChan:    make(chan *server.PlayAudioInfo, _CHANNEL_SIZE),
		skr:         NewSeek(w),
		Status:      0,
		StatusLabel: "Stop",
	}

	p.self.Title = "Status"
	p.self.PaddingLeft = 3

	p.progress.Label = " "
	p.countDown.Text = " 00:00/00:00"

	//p.window.setPersentRect(p.self, 0.14, 0.60, 0.315, 0.27) // status
	p.countDown.SetRect(w.MaxX-_COUNT_DOWN_WIDTH, w.MaxY-_COUNT_DOWN_HEIGHT, w.MaxX, w.MaxY)
	maxX, maxY := p.window.GetMax()
	p.progress.SetRect(0, w.MaxY-_COUNT_DOWN_HEIGHT, w.MaxX-_COUNT_DOWN_WIDTH, w.MaxY)
	p.self.SetRect(int(_VOLUME_WIDTH*maxX), int(maxY-_COUNT_DOWN_HEIGHT-_PLAY_STATUS_HEIGHT),
		int(maxX*(_PLAY_STATUS_WIDTH+_VOLUME_WIDTH)), int(maxY-_COUNT_DOWN_HEIGHT))
	return p
}

func (p *playStatus) InitPrint(info *server.PlayAudioInfo) {
	p.init(info)
	p.print()
}

func (p *playStatus) Cronjob() {
	ticker := time.NewTicker(time.Millisecond * 100)
	for {
		select {
		case <-ticker.C:
			if p.skr.Shielding {
				break
			}
			if p.Status == player.PLAY {
				now := time.Now().Unix()
				if now >= p.Endline {
					p.CurrPro = p.Duration
				} else {
					p.CurrPro = p.Duration - (p.Endline - now)
				}

				p.updateProgress()
				p.printPro()
			}
		case info := <-p.infoChan:
			if p.skr.Shielding {
				break
			}
			p.init(info)
			p.print()
		}
	}
}

func (p *playStatus) SeekAudioFile(second int64) {
	if second >= p.Duration-p.CurrPro { // forward limit
		second = p.Duration - p.CurrPro
	} else if -second > p.CurrPro { // rewind limit
		second = -p.CurrPro
	}

	p.CurrPro += second
	p.Endline -= second
	p.skr.Handle(p.CurrPro, p.AudioPath, p.Status == player.PLAY)
	p.updateProgress()
	p.window.SyncPrint(p.printPro)
}

func (p *playStatus) Notify(info *server.PlayAudioInfo) {
	if info != nil {
		p.infoChan <- info
	}
}

func (p *playStatus) init(info *server.PlayAudioInfo) {
	p.Status = info.Status
	p.LoopMode = guider.Mode(info.Mode)
	p.AudioPath = info.Pathinfo
	p.Duration = info.Duration
	p.CurrPro = info.Curr
	p.Endline = int64(info.Duration-info.Curr) + time.Now().Unix()
	p.updateProgress()
}

func (p *playStatus) print() {
	p.printPro()
	p.printStatus()
}

func (p *playStatus) printStatus() {
	pg := p.self

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

	pg.Text = p.text()
	termui.Render(pg)
}

func (p *playStatus) printPro() {
	termui.Render(p.countDown)
	termui.Render(p.progress)
}

func (p *playStatus) updateProgress() {
	if p.Duration == 0 {
		p.progress.Percent = 0
	} else {
		p.progress.Percent = int(100 * p.CurrPro / p.Duration)
	}
	p.countDown.Text = fmt.Sprintf(" %s/%s", formatDuration(p.CurrPro), formatDuration(p.Duration))
}

func (p *playStatus) text() string {
	return fmt.Sprintf(
		"\n%s\n"+
			"\nAudio:      %s\n"+
			"\nDuration:   %s\n"+
			"\nCycelMode:  %s\n",
		p.StatusLabel, filepath.Base(p.AudioPath), formatDuration(p.Duration), p.formatLoopMode())
}

func (p *playStatus) ChangeLoopMode() {
	newModel := guider.SEQ
	switch p.LoopMode {
	case guider.SEQ:
		newModel = guider.RANDOM
	case guider.RANDOM:
		newModel = guider.SINGLE
	case guider.SINGLE:
		newModel = guider.SEQ
	}
	rpc_client.ChangeLoopModel(newModel)

	info := rpc_client.PlayStatus()
	p.Notify(info)
}

func formatDuration(t int64) string {
	return fmt.Sprintf("%02d:%02d", t/60, t%60)
}

func (p *playStatus) formatLoopMode() string {
	switch p.LoopMode {
	case guider.SINGLE:
		return "Single 🔂"
	case guider.RANDOM:
		return "Random 🔀"
	case guider.SEQ:
		return "Order  🔁"
	default:
		return "Unknow"
	}
}
