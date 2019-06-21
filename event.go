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

// Event is the Event interface.
type Event interface {
}

// EventScreen represents screen changes such as terminal resizing.
type EventScreen struct {
}

// EventKey represents a keypress event.
type EventKey struct {
	Key   Key
	Rune  rune
	State bool
	Shift bool
	Ctrl  bool
	Alt   bool
	Meta  bool
}

func (e EventKey) HasModifiers() bool {
	return e.Shift || e.Ctrl || e.Alt || e.Meta
}

// EventMouse represents a mouse press event.
type EventMouse struct {
	X, Y   int
	Button int
	State  bool
}

// EventQuit represents a quit event.
type EventQuit struct {
}
