package ui

import (
	"github.com/gizak/termui/v3"
	"tap/rpc_client"
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
	maxX, _ := s.window.GetMax()
	s.self.SetRect(int(maxX*(_VOLUME_WIDTH+_PLAY_STATUS_WIDTH)), 0,
		s.window.MaxX, _SEARCH_HEIGHT)
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
	termui.Render(s.self)
}

func (s *searchInput) WidgetKeys() string {
	return "<C-l>   Clear input\n" +
		"<Enter> Entry audio list\n"
}

func (s *searchInput) HandleEvent(input string) {
	switch input {
	case "<C-l>":
		s.self.Reset()
		s.flushAl()
	case "<Enter>":
		s.window.ChoseItem(s.window.audioList)
	default:
		s.self.HandleInput(input)
		s.flushAl()
	}
}

func (s *searchInput) flushAl() {
	rows := rpc_client.Search(s.self.String())
	s.window.audioList.NotifyRowsChange(rows)
}
