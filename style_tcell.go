// +build !disableTCell

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
	"image/color"

	"github.com/gdamore/tcell"
)

// StyleToTCellStyle converts a provided Style into a tcell.Style.
func StyleToTCellStyle(style Style) tcell.Style {
	tStyle := tcell.StyleDefault
	if style.Foreground != ColorNone {
		tStyle = tStyle.Foreground(RGBAToTCellColor(style.Foreground))
	} else {
		tStyle = tStyle.Foreground(-1)
	}
	if style.Background != ColorNone {
		tStyle = tStyle.Background(RGBAToTCellColor(style.Background))
	} else {
		tStyle = tStyle.Background(-1)
	}
	// A little heavy...
	tStyle = tStyle.Bold(style.Bold).Underline(style.Underline).Blink(style.Blink).Dim(style.Dim).Reverse(style.Reverse)

	return tStyle
}

// RGBAToTCellColor converts a color.RGBA to a tcell.Color type.
func RGBAToTCellColor(color color.RGBA) tcell.Color {
	return tcell.NewRGBColor(int32(color.R), int32(color.G), int32(color.B)) & tcell.ColorIsRGB
}
