package glog

import (
	"os"
	"sync"
)

type File struct {
	file *os.File
	sync.Mutex
}

func (f *File) Write(p []byte) (n int, err error) {
	f.Lock()
	defer f.Unlock()
	return f.file.Write(p)
}
