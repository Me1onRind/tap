package local

import (
	"errors"
	"io/ioutil"
	"path/filepath"
	"regexp"
)

type LocalProvider struct {
	dirs  []string
	files map[string][]string

	currDir string
}

func NewLocalProvider(dirs []string) *LocalProvider {
	p := &LocalProvider{
		files: make(map[string][]string),
	}
	for _, v := range dirs {
		if dir, err := filepath.Abs(v); err == nil {
			p.dirs = append(p.dirs, dir)
		}
	}
	currDir, _ := filepath.Abs(dirs[0])
	p.currDir = currDir
	return p
}

func (p *LocalProvider) ListAll() ([]string, error) {
	if items, ok := p.files[p.currDir]; ok {
		return items, nil
	}

	rd, err := ioutil.ReadDir(p.currDir)
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
		files = append(files, p.currDir+"/"+fi.Name())
	}
	p.files[p.currDir] = files
	return files, nil
}

func (p *LocalProvider) ListDirs() []string {
	var ret []string
	for _, v := range p.dirs {
		dir, _ := filepath.Abs(v)
		if len(dir) > 0 {
			ret = append(ret, dir)
		}
	}
	return ret
}

func (p *LocalProvider) CurrDir() string {
	return p.currDir
}

func (p *LocalProvider) Search(reg string) ([]string, error) {
	var files []string
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
	var ret []string

	for _, v := range files {
		if regex.MatchString(filepath.Base(v)) {
			ret = append(ret, v)
		}
	}
	return ret, nil
}

func (p *LocalProvider) Filepath(name string) string {
	return p.currDir + "/" + name
}

func (p *LocalProvider) SetDir(dir string) error {
	for _, v := range p.dirs {
		if v == dir {
			p.currDir = dir
			return nil
		}
	}
	return errors.New("Illegal directory")
}

func (p *LocalProvider) Flush() {
}
