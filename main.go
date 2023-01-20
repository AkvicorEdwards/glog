package glog

import (
	"log"
	"os"
)

const (
	ZERO = 0
	INFO = 1 << iota
	DEBUG
	WARNING
)

var mask = INFO | WARNING

var log_file *os.File
var file_logger *log.Logger

func SetMask(m int) {
	mask = m
}

func SetLogFile(path string) error {
	var err error
	log_file, err = os.OpenFile(path, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	file_logger = log.New(log_file, "", log.LstdFlags|log.Lshortfile)
	return nil
}

func CloseFile() {
	file_logger = nil
	err := log_file.Close()
	if err != nil {
		log.Println("failed to close log file")
	}
	log_file = nil
}

func Info(format string, values ...interface{}) {
	if (mask & INFO) != 0 {
		log.Printf("[INFO] "+format, values...)
		if file_logger != nil {
			file_logger.Printf("[INFO] "+format, values...)
		}
	}
}

func Debug(format string, values ...interface{}) {
	if (mask & DEBUG) != 0 {
		log.Printf("[DEBUG] "+format, values...)
		if file_logger != nil {
			file_logger.Printf("[DEBUG] "+format, values...)
		}
	}
}

func Warning(format string, values ...interface{}) {
	if (mask & WARNING) != 0 {
		log.Printf("[WARNING] "+format, values...)
		if file_logger != nil {
			file_logger.Printf("[WARNING] "+format, values...)
		}
	}
}

func Fatal(format string, values ...interface{}) {
	if file_logger != nil {
		log.Printf("[FATAL] "+format, values...)
		file_logger.Fatalf("[FATAL] "+format, values...)
	} else {
		log.Fatalf("[FATAL] "+format, values...)
	}
}
