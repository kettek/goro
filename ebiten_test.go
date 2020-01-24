// +build enableEbiten,!disableEbiten

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
	"log"

	"testing"
)

func TestMain(m *testing.M) {
	if err := InitEbiten(); err != nil {
		log.Fatal(err)
	}

	Setup(func(screen *Screen) {
		screen.SetTitle("Ebiten Test")
		screen.SetSize(30, 30)

		screen.SetDefaultForeground(ColorWhite)
		screen.SetDefaultBackground(ColorBlack)
	})

	Run(func(screen *Screen) {
		mapScreen, err := NewScreen(30, 30)
		if err != nil {
			log.Fatal(err)
		}

		for {
			mapScreen.DrawString(1, 1, "test", Style{})

			screen.DrawScreen(0, 0, 30, 30, 0, 0, mapScreen)
			screen.Flush()

			mapScreen.DrawString(1, 1, "    ", Style{})

			switch event := screen.WaitEvent().(type) {
			case EventKey:
				if event.Key == KeyEscape {
					Quit()
				}
			case EventQuit:
				return
			}
		}
	})
}
