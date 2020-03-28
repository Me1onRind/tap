package ui

import (
	"github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

type helpBox struct {
	self   *widgets.Paragraph
	window *Window
}

func (h *helpBox) Entry() {
	h.self.BorderStyle.Fg = termui.ColorGreen
	h.Print()
}

func (h *helpBox) Leave() {
	h.self.BorderStyle.Fg = termui.ColorWhite
	h.Print()
}

func (h *helpBox) HandleEvent(intput string) {
	switch intput {
	case "j", "<Down>":
	case "k", "<Up>":
	}
}

func newHelpBox(window *Window) *helpBox {
	h := &helpBox{
		self:   widgets.NewParagraph(),
		window: window,
	}
	h.self.Title = "Global Keys"
	h.self.PaddingLeft = 1
	h.self.Text =
		"<Tab>   Next item\n" +
			"<C-f>   Search\n" +
			"<C-a>   Audio List\n" +
			"<C-k>   Volume up\n" +
			"<C-j>   Volume down\n" +
			"<C-n>   Change Loop Mode\n" +
			"<C-n>   Change Loop Mode\n" +
			"<Left>  Rewind\n" +
			"<Right> Forward\n"
	h.window.setPersentRect(h.self, 0.07, 0.05, 0.2, 0.35)
	return h
}

func (h *helpBox) Print() {
	termui.Render(h.self)
}
