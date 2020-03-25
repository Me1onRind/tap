package ui

import (
	"fmt"
	"github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

type volumeController struct {
	self   *widgets.BarChart // play status
	window *Window
	volume float32
}

func newVolumeController(window *Window, volume float32) *volumeController {
	v := &volumeController{
		self:   widgets.NewBarChart(),
		window: window,
		volume: volume,
	}
	v.self.Data = []float64{float64(volume)}
	v.self.Labels = []string{"volume"}
	v.self.BarGap = 0
	v.self.PaddingTop = -1
	v.self.NumFormatter = func(n float64) string {
		return fmt.Sprintf("%d", int((n+0.0001)*100))
	}
	return v
}

func (v *volumeController) print() {
	v.self.MaxVal = 1.0
	v.self.BarWidth = int(v.window.MaxX * 0.06)
	v.self.Data[0] = float64(v.volume)
	v.window.setPersentRect(v.self, 0.07, 0.61, 0.07, 0.27)
	v.window.syncPrint(v.self)
}

func (v *volumeController) handleEvent(e *termui.Event) {
}

func (v *volumeController) up() {
	if v.volume > 0.99 {
		return
	}
	v.volume += 0.01
	v.window.setVolume(v.volume)
	v.print()
}

func (v *volumeController) down() {
	if v.volume < 0.01 {
		return
	}
	v.volume -= 0.01
	v.window.setVolume(v.volume)
	v.print()
}
