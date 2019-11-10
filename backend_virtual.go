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

import (
	"github.com/kettek/goro/glyphs"
)

// BackendVirtual is the virtual backend used for sub-screens.
type BackendVirtual struct {
	screen      Screen
	refreshChan chan struct{}
	hasStarted  bool
	title       string
}

// Init sets up our appropriate data structures.
func (backend *BackendVirtual) Init() error {
	backend.refreshChan = make(chan struct{})

	backend.title = "goro - Virtual"

	return nil
}

// Quit closes all screens.
func (backend *BackendVirtual) Quit() {
	backend.screen.Close()
}

// Refresh forces the screen to redraw.
func (backend *BackendVirtual) Refresh() {
}

// Setup runs the given function cb.
func (backend *BackendVirtual) Setup(cb func(*Screen)) (err error) {
	return err
}

// Run runs the given function cb as a goroutine.
func (backend *BackendVirtual) Run(cb func(*Screen)) (err error) {
	return err
}

// Size returns the current backend window dimensions.
func (backend *BackendVirtual) Size() (int, int) {
	return 0, 0
}

// SetSize sets the backend window to the provided width and height.
func (backend *BackendVirtual) SetSize(w, h int) {
}

// Units returns the unit type the backend uses for Size().
func (backend *BackendVirtual) Units() int {
	return UnitCells
}

// Scale returns the current backend window scaling.
func (backend *BackendVirtual) Scale() float64 {
	return 1
}

// SetScale sets the backend window's scaling.
func (backend *BackendVirtual) SetScale(scale float64) {
}

// SetTitle sets the backend window's title.
func (backend *BackendVirtual) SetTitle(title string) {
	backend.title = title
}

// SetGlyphs does nothing!
func (backend *BackendVirtual) SetGlyphs(id glyphs.ID, path string, size float64) error {
	return nil
}

// SyncSize does nothing!
func (backend *BackendVirtual) SyncSize() {
	return
}
