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
	"github.com/kettek/goro/glyphs"
)

// Backend is an interface through which a Screen is displayed and controlled.
type Backend interface {
	Init() error
	Run(func(*Screen)) error
	Refresh()
	Quit()
	Size() (int, int)
	SetSize(int, int)
	Units() int
	Scale() float64
	SetScale(float64)
	SetTitle(string)
	SetGlyphs(glyphs.ID, string, float64) error
	SyncSize()
}
