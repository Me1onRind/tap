package local

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"tap/backend"
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

func (p *LocalProvider) Search(reg string) ([]backend.AudioItem, error) {
	log.Println(reg)
	var files []backend.AudioItem
	var err error
	files = p.files[p.currDir]
	if files == nil {
		files, err = p.ListAll()
		if err != nil {
			return nil, err
		}
	}

	regex, err := regexp.Compile(reg)
	if err != nil {
		return nil, err
	}
	var ret []backend.AudioItem

	for _, v := range files {
		fi := v.(os.FileInfo)
		if regex.MatchString(fi.Name()) {
			ret = append(ret, v)
		}
	}
	return ret, nil
}

func (p *LocalProvider) Filepath(name string) (string, error) {
	return filepath.Abs(p.currDir + "/" + name)
}

func (p *LocalProvider) Flush() {
}
