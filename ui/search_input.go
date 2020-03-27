package ui

import (
	"github.com/gizak/termui/v3"
	"tap/ui/custom_widgets"
)

type searchInput struct {
	self   *custom_widgets.Input
	window *Window
}

func newSearchInput(window *Window) *searchInput {
	s := &searchInput{
		self:   custom_widgets.NewInput(),
		window: window,
	}
	s.self.Title = "Search"
	s.self.PaddingLeft = 1
	return s
}

func (s *searchInput) Entry() {
	s.self.BorderStyle.Fg = termui.ColorGreen
	s.self.Force = true
	termui.Render(s.self)
}

func (s *searchInput) Leave() {
	s.self.BorderStyle.Fg = termui.ColorWhite
	s.self.Force = false
	termui.Render(s.self)
}

func (s *searchInput) Print() {
	s.window.setPersentRectWithFixedHeight(s.self, 0.46, 0.05, 0.4, 3)
	termui.Render(s.self)
}

func (s *searchInput) HandleEvent(input string) {
	switch input {
	case "<Enter>":
		s.window.nextItem()
	case "<C-l>":
		s.self.Reset()
		s.flushAl()
	default:
		s.self.HandleInput(input)
		s.flushAl()
	}
}

func (s *searchInput) flushAl() {
	rows := s.window.Search(s.self.String())
	s.window.al.NotifyRowsChange(rows)
}
