// +build windows,cgo,!disableTCell

package goro
/*
This file is a part of goRo, a library for writing roguelikes.
Copyright (C) 2019 Ketchetwahmeegwun T. Southall

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Lesser General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Lesser General Public License for more details.

You should have received a copy of the GNU Lesser General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

// #include <stdlib.h>
import "C"
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
func (backend *BackendTCell) SetTitle(title string) {
	// This is a bogus way to check if we're in a Windows console or not, but...
	if backend.tcellScreen.CharacterSet() == "UTF-16LE" {
		cstr := C.CString(title)
		defer C.free(unsafe.Pointer(cstr))
		_, _, err := procSetConsoleTitleA.Call(
			uintptr(unsafe.Pointer(cstr)),
		)
		if err != syscall.Errno(0) {
			return
		}
	} else {
		fmt.Printf("\033]0;" + title + "\007")
	}
	return
}
