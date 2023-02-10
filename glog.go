package glog

import (
	"fmt"
	"io"
	"os"
	"path"
	"runtime"
	"sync"
	"time"
)

const (
	MaskUNKNOWN = 1 << iota
	MaskDEBUG
	MaskTRACE
	MaskINFO
	MaskWARNING
	MaskERROR
	MaskFATAL
)

const (
	FlagDate = 1 << iota
	FlagTime
	FlagLongFile
	FlagShortFile
	FlagPrefix
	FlagStdFlag = FlagDate | FlagTime
)

var (
	prefixUNKNOWN = "[UNKNOWN] "
	prefixDEBUG   = "[DEBUG  ] "
	prefixTRACE   = "[TRACE  ] "
	prefixINFO    = "[INFO   ] "
	prefixWARNING = "[WARNING] "
	prefixERROR   = "[ERROR  ] "
	prefixFATAL   = "[FATAL  ] "
)

var mask = MaskUNKNOWN | MaskDEBUG | MaskTRACE | MaskINFO | MaskWARNING | MaskERROR | MaskFATAL
var flag = FlagStdFlag

var consoleStdout *File
var consoleStderr *File
var file *File
var lock sync.RWMutex

func init() {
	consoleStdout = new(File)
	consoleStdout.file = os.Stdout
	consoleStderr = new(File)
	consoleStderr.file = os.Stderr
	file = nil
	lock = sync.RWMutex{}
}

func SetMask(m int) {
	lock.Lock()
	defer lock.Unlock()
	mask = m
}

func SetFlag(f int) {
	lock.Lock()
	defer lock.Unlock()
	flag = f
}

func SetLogFile(path string) error {
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	lock.Lock()
	defer lock.Unlock()
	file = new(File)
	file.file = f
	return nil
}

func CloseFile() {
	lock.Lock()
	defer lock.Unlock()
	err := file.file.Close()
	if err != nil {
		write(os.Stderr, prefixERROR, "failed to close log file")
	}
	file = nil
}

func write(w io.Writer, prefix string, format string, values ...any) {
	now := ""
	if (flag&FlagDate) != 0 && (flag&FlagTime) != 0 {
		now = time.Now().Format("2006-01-02 15:04:05 ")
	} else if (flag & FlagDate) != 0 {
		now = time.Now().Format("2006-01-02 ")
	} else if (flag & FlagTime) != 0 {
		now = time.Now().Format("15:04:05 ")
	}
	if (flag & FlagPrefix) != 0 {
		_, _ = w.Write([]byte(now))
		_, _ = w.Write([]byte(prefix))
	} else {
		_, _ = w.Write([]byte(prefix))
		_, _ = w.Write([]byte(now))
	}

	if (flag & FlagLongFile) != 0 {
		_, p, l, ok := runtime.Caller(2)
		if ok {
			_, _ = fmt.Fprintf(w, "%s:%d ", p, l)
		} else {
			_, _ = w.Write([]byte("?:? "))
		}
	} else if (flag & FlagShortFile) != 0 {
		_, p, l, ok := runtime.Caller(2)
		if ok {
			_, _ = fmt.Fprintf(w, "%s:%d ", path.Base(p), l)
		} else {
			_, _ = w.Write([]byte("?:? "))
		}
	}
	_, _ = w.Write([]byte("| "))
	_, _ = fmt.Fprintf(w, format, values...)
	_, _ = w.Write([]byte("\n"))
}

func Unknown(format string, values ...any) {
	lock.Lock()
	defer lock.Unlock()
	if (mask & MaskUNKNOWN) != 0 {
		if file != nil {
			write(file, prefixUNKNOWN, format, values...)
		}
		write(consoleStdout, prefixUNKNOWN, format, values...)
	}
}

func Debug(format string, values ...any) {
	if (mask & MaskDEBUG) != 0 {
		if file != nil {
			write(file, prefixDEBUG, format, values...)
		}
		write(consoleStdout, prefixDEBUG, format, values...)
	}
}

func Trace(format string, values ...any) {
	if (mask & MaskTRACE) != 0 {
		if file != nil {
			write(file, prefixTRACE, format, values...)
		}
		write(consoleStdout, prefixTRACE, format, values...)
	}
}

func Info(format string, values ...any) {
	if (mask & MaskINFO) != 0 {
		if file != nil {
			write(file, prefixINFO, format, values...)
		}
		write(consoleStdout, prefixINFO, format, values...)
	}
}

func Warning(format string, values ...any) {
	if (mask & MaskWARNING) != 0 {
		if file != nil {
			write(file, prefixWARNING, format, values...)
		}
		write(consoleStdout, prefixWARNING, format, values...)
	}
}

func Error(format string, values ...any) {
	if (mask & MaskERROR) != 0 {
		if file != nil {
			write(file, prefixERROR, format, values...)
		}
		write(consoleStderr, prefixERROR, format, values...)
	}
}

func Fatal(format string, values ...any) {
	if (mask & MaskFATAL) != 0 {
		if file != nil {
			write(file, prefixFATAL, format, values...)
		}
		write(consoleStderr, prefixFATAL, format, values...)
	}
}
