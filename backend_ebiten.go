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
	"fmt"

	"github.com/hajimehoshi/ebiten"
)

// BackendEbiten is our Ebiten backend.
type BackendEbiten struct {
	screen        Screen
	imageBuffer   *ebiten.Image
	width, height int
	hasStarted    bool

	pressedKeys  []bool
	pressedMouse []bool
	refreshChan  chan struct{}
}

// InitEbiten initializes the Ebiten backend for use. Calls BackendEbiten.Init().
func InitEbiten() error {
	return Init(Backend(&BackendEbiten{}))
}

// Init sets up our appropriate data structures.
func (backend *BackendEbiten) Init() error {
	backend.pressedKeys = make([]bool, ebiten.KeyMax+1)
	backend.pressedMouse = make([]bool, ebiten.MouseButtonMiddle+1)

	if err := backend.screen.Init(); err != nil {
		return err
	}

	backend.refreshChan = make(chan struct{})

	backend.width = 320
	backend.height = 240

	return nil
}

// Quit closes all screens.
func (backend *BackendEbiten) Quit() {
	backend.screen.Close()
}

// Refresh forces the screen to redraw.
func (backend *BackendEbiten) Refresh() {
}

// Run runs the given function cb as a goroutine.
func (backend *BackendEbiten) Run(cb func(*Screen)) (err error) {
	err = ebiten.Run(func(imageBuffer *ebiten.Image) (err error) {
		if !backend.hasStarted {
			go func() {
				cb(&backend.screen)
			}()
			backend.hasStarted = true
		}

		// ... Ew.
		for k := ebiten.Key(0); k <= ebiten.KeyMax; k++ {
			if ebiten.IsKeyPressed(k) {
				if !backend.pressedKeys[k] {
					fmt.Printf("Sending %d.\n", int16(k))
					if backend.screen.UseKeys {
						backend.screen.eventChan <- Event(EventKey{Key: Key(k), Rune: rune(k)})
					}
				}
				backend.pressedKeys[k] = true
			} else {
				backend.pressedKeys[k] = false
			}
		}
		// ... Ew x2.
		for m := ebiten.MouseButtonLeft; m <= ebiten.MouseButtonMiddle; m++ {
			if ebiten.IsMouseButtonPressed(m) {
				if !backend.pressedMouse[m] && backend.screen.UseMouse {
					backend.screen.eventChan <- Event(EventMouse{})
				}
				backend.pressedMouse[m] = true
			} else {
				backend.pressedMouse[m] = false
			}
		}

		// Draw
		if backend.screen.Redraw && !ebiten.IsDrawingSkipped() {
			for y := 0; y < len(backend.screen.cells); y++ {
				for x := 0; x < len(backend.screen.cells[y]); x++ {
					if backend.screen.cells[y][x].Dirty {
						// text.Draw(backend.screen, backend.screen.cells[y][x].Rune, backend.screen.Font, x*backend.screen.FontWidth(), y*backend.screen.FontHeight(), StyleToEbitenStyle(backend.screen.cells[y][x].Style))
						backend.screen.cells[y][x].Dirty = false
					}
				}
			}
			backend.screen.Redraw = false
		}

		return nil
	}, backend.width, backend.height, 1, "GoingRogue - Ebiten")

	return err
}

// Size returns the current backend window dimensions.
func (backend *BackendEbiten) Size() (int, int) {
	return backend.width, backend.height
}

// SetSize sets the backend window to the provided width and height.
func (backend *BackendEbiten) SetSize(w, h int) {
	backend.width, backend.height = w, h
	ebiten.SetScreenSize(w, h)
}

// Units returns the unit type the backend uses for Size().
func (backend *BackendEbiten) Units() int {
	return UnitPixels
}

// Scale returns the current backend window scaling.
func (backend *BackendEbiten) Scale() float64 {
	return ebiten.ScreenScale()
}

// SetScale sets the backend window's scaling.
func (backend *BackendEbiten) SetScale(scale float64) {
	ebiten.SetScreenScale(scale)
}

// SetTitle sets the backend window's title.
func (backend *BackendEbiten) SetTitle(title string) {
	ebiten.SetWindowTitle(title)
}
