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

package goro

import (
	"os"
)

var globalBackend Backend

// Init initializes a given Backend interface. Built-in options are EbitenBackend and TCellBackend.
func Init(backend Backend) error {
	globalBackend = backend
	return backend.Init()
}

// Quit tells the current Backend to Quit before calling os.Exit.
func Quit() {
	globalBackend.Quit()
	os.Exit(0)
}

// Run runs the Backend with the provided logic callback.
func Run(cb func(*Screen)) error {
	return globalBackend.Run(cb)
}
