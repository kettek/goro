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
	"image/color"
)

// Style represents the styling for a given Cell.
type Style struct {
	Background, Foreground               color.RGBA
	Blink, Underline, Bold, Dim, Reverse bool
}

// Color is an alias to color.RGBA.
type Color = color.RGBA

// These colors represent the default 8 color and 8 bright colors.
var (
	ColorNone    = Color{0x00, 0x00, 0x00, 0x00}
	Color0       = Color{0x00, 0x00, 0x00, 0xFF}
	Color1       = Color{0x80, 0x00, 0x00, 0xFF}
	Color2       = Color{0x00, 0x80, 0x00, 0xFF}
	Color3       = Color{0x80, 0x80, 0x00, 0xFF}
	Color4       = Color{0x00, 0x00, 0x80, 0xFF}
	Color5       = Color{0x80, 0x00, 0x80, 0xFF}
	Color6       = Color{0x00, 0x80, 0x80, 0xFF}
	Color7       = Color{0xC0, 0xC0, 0xC0, 0xFF}
	Color8       = Color{0x80, 0x80, 0x80, 0xFF}
	Color9       = Color{0xFF, 0x00, 0x00, 0xFF}
	Color10      = Color{0x00, 0xFF, 0x00, 0xFF}
	Color11      = Color{0xFF, 0xFF, 0x00, 0xFF}
	Color12      = Color{0x00, 0x00, 0xFF, 0xFF}
	Color13      = Color{0xFF, 0x00, 0xFF, 0xFF}
	Color14      = Color{0x00, 0xFF, 0xFF, 0xFF}
	Color15      = Color{0xFF, 0xFF, 0xFF, 0xFF}
	ColorBlack   = Color0
	ColorMaroon  = Color1
	ColorGreen   = Color2
	ColorOlive   = Color3
	ColorNavy    = Color4
	ColorPurple  = Color5
	ColorTeal    = Color6
	ColorSilver  = Color7
	ColorGray    = Color8
	ColorRed     = Color9
	ColorLime    = Color10
	ColorYellow  = Color11
	ColorBlue    = Color12
	ColorFuchsia = Color13
	ColorAqua    = Color14
	ColorWhite   = Color15
)
