package ui

import (
	"testing"
)

func Test_Window_One(t *testing.T) {
	w := NewWindow(nil)
	w.Init()
	defer w.Close()
}
