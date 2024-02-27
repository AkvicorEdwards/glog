package glog

import (
	"fmt"
	"testing"
)

func TestGlog(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("expected os.Exit(0)")
		} else {
			if fmt.Sprint(r) != "unexpected call to os.Exit(0) during test" {
				t.Errorf("expected os.Exit(0) got %v", r)
			}
		}
	}()
	SetMask(MaskAll)
	SetFlag(FlagAll)
	Unknown("test %s", "Unknown")
	Debug("test %s", "Debug")
	Trace("test %s", "Trace")
	Info("test %s", "Info")
	Warning("test %s", "Warning")
	Error("test %s", "Error")
	fmt.Println()
	SetFlag(FlagStd | FlagShortFile)
	Unknown("test %s", "Unknown")
	Debug("test %s", "Debug")
	Trace("test %s", "Trace")
	Info("test %s", "Info")
	Warning("test %s", "Warning")
	Error("test %s", "Error")
	fmt.Println()
	SetFlag(FlagStd | FlagShortFile | FlagSuffix)
	Unknown("test %s", "Unknown")
	Debug("test %s", "Debug")
	Trace("test %s", "Trace")
	Info("test %s", "Info")
	Warning("test %s", "Warning")
	Error("test %s", "Error")
	fmt.Println()
	SetSeparatorStart(" : ")
	SetSeparatorEnd(" [ ")
	SetSeparatorEndEnd(" ]")
	Unknown("test %s", "Unknown")
	Debug("test %s", "Debug")
	Trace("test %s", "Trace")
	Info("test %s", "Info")
	Warning("test %s", "Warning")
	Error("test %s", "Error")
	fmt.Println()
	SetFlag(FlagAll ^ FlagSuffix)
	Unknown("test %s", "Unknown")
	Debug("test %s", "Debug")
	Trace("test %s", "Trace")
	Info("test %s", "Info")
	Warning("test %s", "Warning")
	Error("test %s", "Error")
	Fatal("test %s", "Fatal")
}
