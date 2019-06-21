// +build windows,!disableTCell

package goro

import (
	"fmt"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	kernel32DLL          = windows.NewLazyDLL("kernel32.dll")
	procSetConsoleTitleA = kernel32DLL.NewProc("SetConsoleTitleA")
)

// SetTitle sets the backend window's title.
func (backend *BackendTCell) SetTitle(title string) error {
	// This is a bogus way to check if we're in a Windows console or not, but...
	if backend.tcellScreen.CharacterSet() == "UTF-16LE" {
		_, _, err := procSetConsoleTitleA.Call(
			uintptr(unsafe.Pointer(StringToCharPtr(title))),
		)
		if err != syscall.Errno(0) {
			return err
		}
	} else {
		fmt.Printf("\033]0;" + title + "\007")
	}
	return nil
}
