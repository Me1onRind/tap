package local

import (
	"errors"
	"tap/backend"
	"io/ioutil"
	"os"
	"path/filepath"
)

type LocalProvider struct {
	dirs  []string
	files map[string][]backend.AudioItem

	currDir string
}

func NewLocalProvider(dirs []string) *LocalProvider {
	p := &LocalProvider{
		dirs:    dirs,
		files:   make(map[string][]backend.AudioItem),
		currDir: dirs[0],
	}

	return p
}

func (p *LocalProvider) ListAll() ([]backend.AudioItem, error) {
	if items, ok := p.files[p.currDir]; ok {
		return items, nil
	}

	rd, err := ioutil.ReadDir(p.currDir)
	if err != nil {
		return nil, err
	}

	var files []backend.AudioItem
	for _, fi := range rd {
		if fi.IsDir() {
			continue
		}
		ext := filepath.Ext(fi.Name())
		if ext != ".mp3" && ext != ".wmv" {
			continue
		}
		files = append(files, fi)
	}
	p.files[p.currDir] = files
	return files, nil
}

func (p *LocalProvider) Search(input string) ([]backend.AudioItem, error) {
	return p.ListAll()
}

func (p *LocalProvider) Filepath(index int) (string, error) {
	items := p.files[p.currDir]
	if index >= len(items) {
		return "", errors.New("Out of range")
	}
	fi := items[index].(os.FileInfo)
	return filepath.Abs(p.currDir + "/" + fi.Name())
}

func (p *LocalProvider) Flush() {
}
