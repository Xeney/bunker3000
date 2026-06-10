package colors

import (
	"runtime"
	"syscall"
	"unsafe"
)

var enabled bool

const (
	Reset = "\033[0m"

	Bold  = "\033[1m"
	Under = "\033[4m"

	Red     = "\033[31m"
	Green   = "\033[32m"
	Yellow  = "\033[33m"
	Blue    = "\033[34m"
	Magenta = "\033[35m"
	Cyan    = "\033[36m"

	BrightRed     = "\033[91m"
	BrightGreen   = "\033[92m"
	BrightYellow  = "\033[93m"
	BrightCyan    = "\033[96m"
	BrightWhite   = "\033[97m"
)

func S(text string, color string) string {
	if !enabled {
		return text
	}
	return color + text + Reset
}

func Init() {
	if runtime.GOOS == "windows" {
		enabled = enableVT()
	} else {
		enabled = true
	}
}

func enableVT() bool {
	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	setMode := kernel32.NewProc("SetConsoleMode")
	getMode := kernel32.NewProc("GetConsoleMode")

	handle, err := syscall.GetStdHandle(syscall.STD_OUTPUT_HANDLE)
	if err != nil {
		return false
	}

	var mode uint32
	ret, _, _ := getMode.Call(uintptr(handle), uintptr(unsafe.Pointer(&mode)))
	if ret == 0 {
		return false
	}

	mode |= 0x0004
	ret, _, _ = setMode.Call(uintptr(handle), uintptr(mode))
	return ret != 0
}
