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
	"errors"
	"sync"

	"github.com/kettek/goro/glyphs"
)

// Screen is a virtual Rows x Columns buffer used for drawing runes to.
type Screen struct {
	ScrollX, ScrollY int
	Columns, Rows    int
	cells            [][]Cell
	active           bool
	eventChan        chan Event
	UseKeys          bool
	UseMouse         bool
	AutoSize         bool
	Foreground       Color
	Background       Color
	Redraw           bool
	cellsMutex       sync.Mutex
}

// Init initializes the Screen's data structures and default values.
func (screen *Screen) Init() (err error) {
	screen.eventChan = make(chan Event, 10)
	// Assign our Screen size to either the backend's columns and rows if it uses Cells for units, otherwise use a standard 80x24 size.
	if globalBackend != nil {
		if globalBackend.Units() == UnitCells {
			screen.Columns, screen.Rows = globalBackend.Size()
		} else {
			screen.Columns, screen.Rows = 80, 24
		}
	}
	screen.active = true
	screen.UseKeys = true
	screen.AutoSize = true

	return screen.Sync()
}

// WaitEvent returns an Event from the Screen's event channel.
func (screen *Screen) WaitEvent() Event {
	return <-screen.eventChan
}

// Close sets the Screen as inactive thereby allowing the backend to clean it up.
func (screen *Screen) Close() {
	screen.active = false
}

// Sync synchronizes the screen's actual cells with the current Rows and Columns.
func (screen *Screen) Sync() (err error) {
	screen.cellsMutex.Lock()
	defer screen.cellsMutex.Unlock()
	currRows := len(screen.cells)
	// Grow or shrink our rows.
	if currRows < screen.Rows {
		screen.cells = append(screen.cells, make([][]Cell, screen.Rows-currRows)...)
	} else if currRows > screen.Rows {
		screen.cells = screen.cells[:screen.Rows]
	}
	// Iterate through our rows to grow or shrink their columns.
	for y := range screen.cells {
		currColumns := len(screen.cells[y])
		// Grow or shrink our columns.
		if currColumns < screen.Columns {
			screen.cells[y] = append(screen.cells[y], make([]Cell, screen.Columns-currColumns)...)
			// Mark new cells as dirty
			for x := currColumns; x < screen.Columns; x++ {
				screen.cells[y][x].Dirty = true
			}
		} else if currColumns > screen.Columns {
			screen.cells[y] = screen.cells[y][:screen.Columns]
		}
	}
	return nil
}

func (screen *Screen) checkBounds(x, y int) error {
	if y < 0 || y >= len(screen.cells) {
		return errors.New("y out of range")
	}
	if x < 0 || x >= len(screen.cells[y]) {
		return errors.New("x out of range")
	}
	return nil
}

// DrawRune draws a given rune at the position of x and y with a given style.
func (screen *Screen) DrawRune(x int, y int, r rune, s Style) error {
	screen.cellsMutex.Lock()
	defer screen.cellsMutex.Unlock()
	if err := screen.checkBounds(x, y); err != nil {
		return err
	}
	if screen.cells[y][x].PendingRune == r && screen.cells[y][x].PendingStyle == s {
		return nil
	}
	screen.cells[y][x].PendingRune = r
	screen.cells[y][x].PendingStyle = s
	screen.cells[y][x].Dirty = true
	return nil
}

// DrawString draws a string at the position of x and y with a given style, iterating in the x direction as it goes.
func (screen *Screen) DrawString(x int, y int, str string, s Style) error {
	origX := x
	for _, r := range str {
		if r == '\n' {
			x = origX
			y++
		} else {
			if err := screen.DrawRune(x, y, r, s); err != nil {
				return err
			}
			x++
		}
	}
	return nil
}

// SetForeground sets the foreground at the given location.
func (screen *Screen) SetForeground(x int, y int, c Color) error {
	screen.cellsMutex.Lock()
	defer screen.cellsMutex.Unlock()
	if err := screen.checkBounds(x, y); err != nil {
		return err
	}
	if screen.cells[y][x].PendingStyle.Foreground == c {
		return nil
	}
	screen.cells[y][x].PendingStyle.Foreground = c
	screen.cells[y][x].Dirty = true
	return nil
}

// SetBackground sets the background at the given location.
func (screen *Screen) SetBackground(x int, y int, c Color) error {
	screen.cellsMutex.Lock()
	defer screen.cellsMutex.Unlock()
	if err := screen.checkBounds(x, y); err != nil {
		return err
	}
	if screen.cells[y][x].PendingStyle.Background == c {
		return nil
	}
	screen.cells[y][x].PendingStyle.Background = c
	screen.cells[y][x].Dirty = true
	return nil
}

// SetStyle sets the style at the given location.
func (screen *Screen) SetStyle(x int, y int, s Style) error {
	screen.cellsMutex.Lock()
	defer screen.cellsMutex.Unlock()
	if err := screen.checkBounds(x, y); err != nil {
		return err
	}
	if screen.cells[y][x].PendingStyle == s {
		return nil
	}
	screen.cells[y][x].PendingStyle = s
	screen.cells[y][x].Dirty = true
	return nil
}

// SetRune sets the rune at the given location.
func (screen *Screen) SetRune(x int, y int, r rune) error {
	screen.cellsMutex.Lock()
	defer screen.cellsMutex.Unlock()
	if err := screen.checkBounds(x, y); err != nil {
		return err
	}
	if screen.cells[y][x].PendingRune == r {
		return nil
	}
	screen.cells[y][x].PendingRune = r
	screen.cells[y][x].Dirty = true
	return nil
}

// SetGlyphsID sets the glyphs at a given location.
func (screen *Screen) SetGlyphsID(x int, y int, id glyphs.ID) error {
	screen.cellsMutex.Lock()
	defer screen.cellsMutex.Unlock()
	if err := screen.checkBounds(x, y); err != nil {
		return err
	}
	if screen.cells[y][x].PendingGlyphs == id {
		return nil
	}
	screen.cells[y][x].PendingGlyphs = id
	screen.cells[y][x].Dirty = true
	return nil
}

// Clear clears the underlying screen.
func (screen *Screen) Clear() {
	for y := 0; y < len(screen.cells); y++ {
		for x := 0; x < len(screen.cells[y]); x++ {
			screen.DrawRune(x, y, ' ', Style{})
		}
	}
	screen.Flush()
}

// Flush forcibly causes the screen to commit any pending changes via a Draw* call and render to the backend.
func (screen *Screen) Flush() {
	screen.cellsMutex.Lock()
	for y := 0; y < len(screen.cells); y++ {
		for x := 0; x < len(screen.cells[y]); x++ {
			if screen.cells[y][x].Dirty {
				screen.cells[y][x].Rune = screen.cells[y][x].PendingRune
				screen.cells[y][x].Style = screen.cells[y][x].PendingStyle
				screen.cells[y][x].Glyphs = screen.cells[y][x].PendingGlyphs
				screen.cells[y][x].Dirty = false
				screen.cells[y][x].Redraw = true
			}
		}
	}
	screen.Redraw = true
	defer screen.cellsMutex.Unlock()
	// hmm. We're calling this here so we can force render the view.
	globalBackend.Refresh()
}

// ForceRedraw marks each cell of the screen to be redrawn.
func (screen *Screen) ForceRedraw() {
	screen.cellsMutex.Lock()
	defer screen.cellsMutex.Unlock()
	for y := 0; y < len(screen.cells); y++ {
		for x := 0; x < len(screen.cells[y]); x++ {
			screen.cells[y][x].Redraw = true
		}
	}
	screen.Redraw = true
}

// SetDefaultForeground sets the default foreground color to use.
func (screen *Screen) SetDefaultForeground(c Color) {
	screen.Foreground = c
	screen.ForceRedraw()
}

// SetDefaultBackground sets the default background color to use.
func (screen *Screen) SetDefaultBackground(c Color) {
	screen.Background = c
	screen.ForceRedraw()
}

// Size returns the current screen columns by rows.
func (screen *Screen) Size() (int, int) {
	return screen.Columns, screen.Rows
}

// SetSize sets the screen's columns and rows.
func (screen *Screen) SetSize(c, r int) {
	screen.Columns = c
	screen.Rows = r
	screen.Sync()
	screen.ForceRedraw()
	globalBackend.SyncSize()
}

// WindowSize returns the current backend's window size in its preferred units, if available.
func (screen *Screen) WindowSize() (int, int) {
	return globalBackend.Size()
}

// SetWindowSize sets the current backend's window to the provided width and height in the backend's units, if available.
func (screen *Screen) SetWindowSize(w, h int) {
	globalBackend.SetSize(w, h)
}

// Scale returns the current screen window scaling. Does nothing.
func (screen *Screen) Scale() float64 {
	return globalBackend.Scale()
}

// SetScale sets the screen window's scaling. Does nothing.
func (screen *Screen) SetScale(scale float64) {
	globalBackend.SetScale(scale)
}

// SetTitle sets the backend window's title to title.
func (screen *Screen) SetTitle(title string) {
	globalBackend.SetTitle(title)
}

// SetGlyphs sets the screen backend's window to use the provided font. Only available for graphical backends.
func (screen *Screen) SetGlyphs(id glyphs.ID, path string, size float64) error {
	return globalBackend.SetGlyphs(id, path, size)
}
