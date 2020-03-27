package ui

import (
	"fmt"
	"github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"tap/server"
)

type volumeController struct {
	self   *widgets.BarChart // play status
	window *Window
	volume float32
}

func newVolumeController(window *Window) *volumeController {
	v := &volumeController{
		self:   widgets.NewBarChart(),
		window: window,
	}
	v.self.Data = []float64{float64(0.5)}
	v.self.Labels = []string{"volume"}
	v.self.BarGap = 0
	v.self.PaddingTop = -1
	v.self.NumFormatter = func(n float64) string {
		return fmt.Sprintf("%d", int((n+0.0001)*100))
	}
	v.self.MaxVal = 1.0
	v.self.BarWidth = int(v.window.MaxX * 0.06)
	v.window.setPersentRect(v.self, 0.07, 0.61, 0.07, 0.27)
	return v
}

func (v *volumeController) InitPrint(info *server.PlayAudioInfo) {
	v.volume = info.Volume
	v.Print()
}

func (v *volumeController) Up() {
	if v.volume > 0.99 {
		return
	}
	v.volume += 0.01
	v.window.SetVolume(v.volume)
	v.Print()
}

func (v *volumeController) Down() {
	if v.volume < 0.01 {
		return
	}
	v.volume -= 0.01
	v.window.SetVolume(v.volume)
	v.Print()
}

func (v *volumeController) Print() {
	v.self.Data[0] = float64(v.volume)
	termui.Render(v.self)
}
