package ui

import (
	"github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

type output struct {
	self   *widgets.Paragraph
	window *Window
}

func newOutput(window *Window) *output {
	o := &output{
		self:   widgets.NewParagraph(),
		window: window,
	}
	o.self.Title = "Output"
	o.self.PaddingLeft = 1
	o.self.Text = "Welcome"
	maxX, _ := o.window.GetMax()
	o.self.SetRect(0, 0,
		int(maxX*(_VOLUME_WIDTH+_PLAY_STATUS_WIDTH)),
		o.window.MaxY-_PLAY_STATUS_HEIGHT-_COUNT_DOWN_HEIGHT-_HELP_BOX_HEIGHT)
	//h.window.setPersentRect(h.self, 0.07, 0.05, 0.2, 0.3)
	return o
}

func (o *output) Print() {
	termui.Render(o.self)
}
