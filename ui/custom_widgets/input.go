package custom_widgets

import (
	"github.com/gizak/termui/v3"
	"image"
)

type Input struct {
	termui.Block
	cursorPoint image.Point

	text    []rune
	textLen int

	Force bool
}

func NewInput() *Input {
	input := &Input{
		Block:       *termui.NewBlock(),
		cursorPoint: image.Pt(0, 0),
		text:        make([]rune, 32),
		textLen:     0,
	}
	return input
}

func (self *Input) Draw(buf *termui.Buffer) {
	self.Block.Draw(buf)

	begin := self.Inner.Min.X
	y := self.Inner.Min.Y
	for offset, i := begin, 0; i < self.textLen; offset, i = offset+1, i+1 {
		ch := self.text[i]
		if self.Force && i == self.cursorPoint.X {
			buf.SetCell(termui.NewCell(ch, termui.NewStyle(termui.ColorWhite, termui.ColorMagenta)),
				image.Pt(offset, y))
		} else {
			buf.SetCell(termui.NewCell(ch), image.Pt(offset, y))
		}
		if len(string(ch)) > 1 {
			offset++
		}
	}
	if self.Force && self.cursorPoint.X >= self.textLen {
		buf.SetCell(termui.NewCell('█'), image.Pt(begin+self.cursorPoint.X, y))
	}
}

func (self *Input) Insert(input string) {
	newLen := self.textLen + len(input)
	if newLen+3 > self.Dx() {
		return
	}
	if newLen > len(self.text) {
		newText := make([]rune, newLen+10)
		copy(newText, self.text)
		self.text = newText
	}

	currX := self.cursorPoint.X

	if currX < self.textLen {
		for i := self.textLen - 1; i >= currX; i-- {
			self.text[i+len(input)] = self.text[i]
		}
	}

	for _, v := range input {
		self.text[self.textLen] = v
		self.textLen++
		if len(string(v)) > 1 {
			self.cursorPoint.X += 2
		} else {
			self.cursorPoint.X++
		}
	}
}

func (self *Input) Backspace() {
	if self.textLen > 0 {
		self.textLen--
		if len(string(self.text[self.textLen])) > 1 {
			self.cursorPoint.X -= 2
		} else {
			self.cursorPoint.X--
		}
	}
}

func (self *Input) MoveCursorLeft() {
	if self.cursorPoint.X > 0 {
		self.cursorPoint.X--
	}
}

func (self *Input) MoveCursorRight() {
	if self.cursorPoint.X < self.textLen {
		self.cursorPoint.X++
	}
}

func (self *Input) HandleInput(input string) {
	switch input {
	case "<Backspace>":
		self.Backspace()
	case "<Space>":
		self.Insert(" ")
	case "<Left>":
		self.MoveCursorLeft()
	case "<C-j>":
		self.MoveCursorLeft()
	case "<Right>":
		self.MoveCursorRight()
	case "<C-k>":
		self.MoveCursorRight()
	default:
		self.Insert(input)
	}
}

func (self *Input) String() string {
	if self.textLen == 0 {
		return ""
	}
	return string(self.text[0:self.textLen])
}

func (self *Input) Reset() {
	self.textLen = 0
	self.cursorPoint.X = 0
}
