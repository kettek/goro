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
	Redraw           bool
}

// Init initializes the Screen's data structures and default values.
func (screen *Screen) Init() (err error) {
	screen.eventChan = make(chan Event, 10)
	screen.Columns = 80
	screen.Rows = 24
	screen.active = true
	screen.UseKeys = true

	screen.cells = make([][]Cell, screen.Rows)
	for y := 0; y < screen.Rows; y++ {
		screen.cells[y] = make([]Cell, screen.Columns)
	}

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
	if len(screen.cells) < screen.Rows {
		// grow Y
		for y := len(screen.cells) - 1; y < screen.Rows; y++ {
			if len(screen.cells[y])-1 < screen.Columns {
				// grow X
			} else if len(screen.cells[y])-1 > screen.Columns {
				// shrink X
			}
		}
	} else if len(screen.cells)-1 > screen.Rows {
		for y := len(screen.cells) - 1; y < screen.Rows; y++ {
			// shrink Y
			if len(screen.cells[y])-1 < screen.Columns {
				// grow X
			} else if len(screen.cells[y])-1 > screen.Columns {
				// shrink X
			}
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
	if err := screen.checkBounds(x, y); err != nil {
		return err
	}
	screen.cells[y][x].PendingRune = r
	screen.cells[y][x].PendingStyle = s
	screen.cells[y][x].Dirty = true
	return nil
}

// DrawString draws a string at the position of x and y with a given style, iterating in the x direction as it goes.
func (screen *Screen) DrawString(x int, y int, str string, s Style) error {
	for _, r := range str {
		if err := screen.DrawRune(x, y, r, s); err != nil {
			return err
		}
		x++
	}
	return nil
}

// SetForeground sets the foreground at the given location.
func (screen *Screen) SetForeground(x int, y int, c Color) error {
	if err := screen.checkBounds(x, y); err != nil {
		return err
	}
	screen.cells[y][x].PendingStyle.Foreground = c
	screen.cells[y][x].Dirty = true
	return nil
}

// SetBackground sets the background at the given location.
func (screen *Screen) SetBackground(x int, y int, c Color) error {
	if err := screen.checkBounds(x, y); err != nil {
		return err
	}
	screen.cells[y][x].PendingStyle.Background = c
	screen.cells[y][x].Dirty = true
	return nil
}

// SetStyle sets the style at the given location.
func (screen *Screen) SetStyle(x int, y int, s Style) error {
	if err := screen.checkBounds(x, y); err != nil {
		return err
	}
	screen.cells[y][x].PendingStyle = s
	screen.cells[y][x].Dirty = true
	return nil
}

// SetRune sets the rune at the given location.
func (screen *Screen) SetRune(x int, y int, r rune) error {
	if err := screen.checkBounds(x, y); err != nil {
		return err
	}
	screen.cells[y][x].PendingRune = r
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
	for y := 0; y < len(screen.cells); y++ {
		for x := 0; x < len(screen.cells[y]); x++ {
			if screen.cells[y][x].Dirty {
				screen.cells[y][x].Rune = screen.cells[y][x].PendingRune
				screen.cells[y][x].Style = screen.cells[y][x].PendingStyle
				screen.cells[y][x].Dirty = false
				screen.cells[y][x].Redraw = true
			}
		}
	}
	screen.Redraw = true
	// hmm. We're calling this here so we can force render the view.
	globalBackend.Refresh()
}

// Size returns the current screen window dimensions.
func (screen *Screen) Size() (int, int) {
	return globalBackend.Size()
}

// SetSize sets the screen window to the provided width and height. Does nothing.
func (screen *Screen) SetSize(w, h int) {
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
