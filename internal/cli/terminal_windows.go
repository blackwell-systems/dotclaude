//go:build windows

package cli

import (
	"os"
	"syscall"
	"unsafe"
)

var (
	kernel32                        = syscall.NewLazyDLL("kernel32.dll")
	procGetConsoleMode              = kernel32.NewProc("GetConsoleMode")
	procSetConsoleMode              = kernel32.NewProc("SetConsoleMode")
	enableVirtualTerminalProcessing = uint32(0x0004)
)

// initTerminalPlatform enables ANSI escape codes on Windows 10+
func initTerminalPlatform() {
	// Get stdout handle
	stdout := syscall.Handle(os.Stdout.Fd())

	var mode uint32
	r, _, _ := procGetConsoleMode.Call(uintptr(stdout), uintptr(unsafe.Pointer(&mode)))
	if r == 0 {
		// Failed to get console mode, disable colors
		disableColors()
		return
	}

	// Try to enable virtual terminal processing
	r, _, _ = procSetConsoleMode.Call(uintptr(stdout), uintptr(mode|enableVirtualTerminalProcessing))
	if r == 0 {
		// Windows 10 version 1511 or earlier, or legacy console
		// Disable colors since ANSI won't work
		disableColors()
	}
}
