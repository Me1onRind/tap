package local

import (
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strings"
)

type LocalProvider struct {
	files map[string][]string
}

func NewLocalProvider() *LocalProvider {
	return &LocalProvider{
		files: make(map[string][]string),
	}
}

func (p *LocalProvider) ListAll(dir string) ([]string, error) {
	if items, ok := p.files[dir]; ok {
		return items, nil
	}

	if !strings.HasSuffix(dir, "/") {
		dir += "/"
	}

	rd, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var files []string
	for _, fi := range rd {
		if fi.IsDir() {
			continue
		}
		ext := filepath.Ext(fi.Name())
		if ext != ".mp3" && ext != ".wmv" {
			continue
		}
		files = append(files, dir+fi.Name())
	}
	p.files[dir] = files
	return files, nil
}

func (p *LocalProvider) Search(reg string, dir string) ([]string, error) {
	files, err := p.ListAll(dir)
	regex, err := regexp.Compile(reg)
	if err != nil {
		return nil, err
	}
	var ret []string

	for _, v := range files {
		if regex.MatchString(filepath.Base(v)) {
			ret = append(ret, v)
		}
	}
	return ret, nil
}

func (p *LocalProvider) Flush() {
}
