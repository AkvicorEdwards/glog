package glog

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"sync/atomic"
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

	MaskStd = MaskINFO | MaskWARNING | MaskERROR | MaskFATAL
	MaskAll = MaskUNKNOWN | MaskDEBUG | MaskTRACE | MaskINFO | MaskWARNING | MaskERROR | MaskFATAL
)

const (
	FlagDate = 1 << iota
	FlagTime
	FlagLongFile
	FlagShortFile
	FlagFunc
	FlagPrefix
	FlagSuffix

	FlagStd = FlagDate | FlagTime | FlagPrefix
	FlagAll = FlagDate | FlagTime | FlagShortFile | FlagFunc | FlagPrefix | FlagSuffix
)

var (
	prefixUNKNOWN = "[UNKNOWN]"
	prefixDEBUG   = "[DEBUG  ]"
	prefixTRACE   = "[TRACE  ]"
	prefixINFO    = "[INFO   ]"
	prefixWARNING = "[WARNING]"
	prefixERROR   = "[ERROR  ]"
	prefixFATAL   = "[FATAL  ]"
)

var separatorStart = " : "
var separatorEnd = "  [ "
var separatorEndEnd = " ]"

var mask uint32 = MaskStd
var flag uint32 = FlagStd

var stdout = NewPaperFromFile(os.Stdout)
var stderr = NewPaperFromFile(os.Stderr)
var file = NewPaperFromFile(nil)

func SetSeparatorStart(sep string) {
	separatorStart = sep
}

func SetSeparatorEnd(sep string) {
	separatorEnd = sep
}

func SetSeparatorEndEnd(end string) {
	separatorEndEnd = end
}

func SetMask(m uint32) {
	atomic.StoreUint32(&mask, m)
}

func GetMask() uint32 {
	return atomic.LoadUint32(&mask)
}

func SetFlag(f uint32) {
	atomic.StoreUint32(&flag, f)
}

func GetFlag() uint32 {
	return atomic.LoadUint32(&flag)
}

func SetLogFile(path string) error {
	f, err := os.OpenFile(filepath.Clean(path), os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	file.SetFile(f)
	return nil
}

func CloseFile() {
	err := file.Close()
	if err != nil {
		Error("failed to close log file: %v", err)
	}
}

func CloseStdout() {
	err := stdout.Close()
	if err != nil {
		Error("failed to close log stdout: %v", err)
	}
}

func CloseStderr() {
	err := stderr.Close()
	if err != nil {
		Error("failed to close log stderr: %v", err)
	}
}

func Close() {
	CloseFile()
	CloseStdout()
	CloseStderr()
}

func write(w *Paper, prefix string, format string, values ...any) {
	if !w.Ready() {
		return
	}
	logPrefix := strings.Builder{}
	logSuffix := strings.Builder{}
	now := ""
	flg := GetFlag()

	if (flg&FlagDate) != 0 && (flg&FlagTime) != 0 {
		now = time.Now().Format("2006-01-02 15:04:05")
	} else if (flg & FlagDate) != 0 {
		now = time.Now().Format("2006-01-02")
	} else if (flg & FlagTime) != 0 {
		now = time.Now().Format("15:04:05")
	}
	if (flg & FlagPrefix) != 0 {
		logPrefix.WriteString(now)
		logPrefix.WriteByte(' ')
		logPrefix.WriteString(prefix)
	} else {
		logPrefix.WriteString(prefix)
		logPrefix.WriteByte(' ')
		logPrefix.WriteString(now)
	}

	if (flg & FlagFunc) != 0 {
		c, _, _, ok := runtime.Caller(2)
		if ok {
			logSuffix.WriteString(runtime.FuncForPC(c).Name())
		} else {
			logSuffix.WriteByte('?')
		}
	}

	if (flg & FlagLongFile) != 0 {
		if logSuffix.Len() != 0 {
			logSuffix.WriteByte(' ')
		}
		_, p, l, ok := runtime.Caller(2)
		if ok {
			logSuffix.WriteString(fmt.Sprintf("%s:%d", p, l))
		} else {
			logSuffix.WriteString("?:?")
		}
	} else if (flg & FlagShortFile) != 0 {
		if logSuffix.Len() != 0 {
			logSuffix.WriteByte(' ')
		}
		_, p, l, ok := runtime.Caller(2)
		if ok {
			logSuffix.WriteString(fmt.Sprintf("%s:%d", path.Base(p), l))
		} else {
			logSuffix.WriteString("?:?")
		}
	}

	if (flag & FlagSuffix) != 0 {
		logPrefix.WriteString(separatorStart)
		logPrefix.WriteString(fmt.Sprintf(format, values...))
		logPrefix.WriteString(separatorEnd)
		logPrefix.WriteString(logSuffix.String())
		logPrefix.WriteString(separatorEndEnd)
	} else {
		logPrefix.WriteByte(' ')
		logPrefix.WriteString(logSuffix.String())
		logPrefix.WriteString(separatorStart)
		logPrefix.WriteString(fmt.Sprintf(format, values...))
	}

	logPrefix.WriteByte('\n')
	w.WriteString(logPrefix.String())
}

func Unknown(format string, values ...any) {
	if (GetMask() & MaskUNKNOWN) != 0 {
		write(file, prefixUNKNOWN, format, values...)
		write(stdout, prefixUNKNOWN, format, values...)
	}
}

func Debug(format string, values ...any) {
	if (GetMask() & MaskDEBUG) != 0 {
		write(file, prefixDEBUG, format, values...)
		write(stdout, prefixDEBUG, format, values...)
	}
}

func Trace(format string, values ...any) {
	if (GetMask() & MaskTRACE) != 0 {
		write(file, prefixTRACE, format, values...)
		write(stdout, prefixTRACE, format, values...)
	}
}

func Info(format string, values ...any) {
	if (GetMask() & MaskINFO) != 0 {
		write(file, prefixINFO, format, values...)
		write(stdout, prefixINFO, format, values...)
	}
}

func Warning(format string, values ...any) {
	if (GetMask() & MaskWARNING) != 0 {
		write(file, prefixWARNING, format, values...)
		write(stdout, prefixWARNING, format, values...)
	}
}

func Error(format string, values ...any) {
	if (GetMask() & MaskERROR) != 0 {
		write(file, prefixERROR, format, values...)
		write(stderr, prefixERROR, format, values...)
	}
}

func Fatal(format string, values ...any) {
	if (GetMask() & MaskFATAL) != 0 {
		write(file, prefixFATAL, format, values...)
		write(stderr, prefixFATAL, format, values...)
		CloseFile()
		os.Exit(0)
	}
}
