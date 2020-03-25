package local

import (
	"testing"
)

func Test_Local_One(t *testing.T) {
	lp := NewSingleDirProvider("./")
	list, _ := lp.AllList()
	for _, v := range list {
		t.Log(lp.Filepath(v))
	}
}
