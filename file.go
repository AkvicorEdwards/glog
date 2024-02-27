package glog

import (
	"os"
	"sync"
)

type Paper struct {
	file *os.File
	sync.Mutex
}

func NewPaperFromFile(file *os.File) *Paper {
	return &Paper{
		file:  file,
		Mutex: sync.Mutex{},
	}
}

func (f *Paper) SetFile(file *os.File) {
	f.Lock()
	defer f.Unlock()
	if f.file != nil {
		_ = f.file.Close()
	}
	if file != nil {
		f.file = file
	}
}

func (f *Paper) WriteString(s string) {
	f.Lock()
	defer f.Unlock()
	if f.file != nil {
		_, _ = f.file.WriteString(s)
	}
}

func (f *Paper) Ready() bool {
	f.Lock()
	defer f.Unlock()
	return f.file != nil
}

func (f *Paper) Close() error {
	f.Lock()
	defer f.Unlock()
	var err error = nil
	if f.file != nil {
		err = f.file.Close()
		f.file = nil
	}
	return err
}
