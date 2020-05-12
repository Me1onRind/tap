package ui

import (
	"fmt"
	"github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"tap/rpc_client"
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
	//v.self.BarWidth =
	maxX, maxY := v.window.GetMax()
	v.self.BarWidth = int(maxX * _VOLUME_WIDTH)
	v.self.SetRect(0, int(maxY-_COUNT_DOWN_HEIGHT-_PLAY_STATUS_HEIGHT),
		int(maxX*_VOLUME_WIDTH), int(maxY-_COUNT_DOWN_HEIGHT))

	//v.window.setPersentRect(v.self, 0.07, 0.60, 0.07, 0.27)
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
	rpc_client.SetVolume(v.volume)
	v.Print()
}

func (v *volumeController) Down() {
	if v.volume < 0.01 {
		return
	}
	v.volume -= 0.01
	rpc_client.SetVolume(v.volume)
	v.Print()
}

func (v *volumeController) Print() {
	v.self.Data[0] = float64(v.volume)
	termui.Render(v.self)
}
