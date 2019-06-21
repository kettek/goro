// +build !disableTCell

package goro

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

func RGBAToTCellColor(color color.RGBA) tcell.Color {
	return tcell.NewRGBColor(int32(color.R), int32(color.G), int32(color.B))
}
