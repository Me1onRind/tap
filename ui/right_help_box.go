package ui

import (
	"github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

type rightHelpBox struct {
	self   *widgets.Paragraph
	window *Window
}

func newRightHelpBox(window *Window) *rightHelpBox {
	r := &rightHelpBox{
		self:   widgets.NewParagraph(),
		window: window,
	}
	r.self.Title = "Widget Keys"
	r.self.PaddingLeft = 1
	r.self.Text =
		"<C-k>   Volume up\n" +
			"<C-j>   Volume down\n" +
			"<C-n>   Change Loop Mode\n" +
			"<Left>  Rewind\n" +
			"<Right> Forward\n" +
			"<Esc>   Exit\n"

	maxX, _ := r.window.GetMax()
	r.self.SetRect(int(maxX*(_VOLUME_WIDTH+_PLAY_STATUS_WIDTH+_LIST_WIDTH)), _SEARCH_HEIGHT,
		int(maxX*(_VOLUME_WIDTH+_PLAY_STATUS_WIDTH+_LIST_WIDTH+_RIGHT_HELP_BOX_WIDTH)),
		r.window.MaxY-_COUNT_DOWN_HEIGHT)
	return r
}

func (r *rightHelpBox) Print() {
	termui.Render(r.self)
}
