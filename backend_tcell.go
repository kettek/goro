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
	"github.com/gdamore/tcell"
)

// BackendTCell is the backend for the tcell library.
type BackendTCell struct {
	screen      Screen
	tcellScreen tcell.Screen
	refreshChan chan struct{}
}

// InitTCell initializes the TCell backend for use. Calls BackendTCell.Init().
func InitTCell() error {
	return Init(Backend(&BackendTCell{}))
}

// Init initializes the base data structures and sets up the library for use.
func (backend *BackendTCell) Init() (err error) {
	backend.tcellScreen, err = tcell.NewScreen()
	if err != nil {
		return err
	}

	if err := backend.tcellScreen.Init(); err != nil {
		return err
	}

	backend.tcellScreen.EnableMouse()
	backend.tcellScreen.SetStyle(tcell.StyleDefault)
	backend.tcellScreen.Clear()

	if err := backend.screen.Init(); err != nil {
		return err
	}

	backend.refreshChan = make(chan struct{})

	return nil
}

// Quit causes our screen to close.
func (backend *BackendTCell) Quit() {
	backend.screen.Close()
	backend.tcellScreen.Fini()
}

// Run runs the given function cb as a goroutine and starts the entire tcell loop.
func (backend *BackendTCell) Run(cb func(*Screen)) (err error) {
	go func() {
		cb(&backend.screen)
	}()

	// I guess this is okay to do.
	go func() {
		for {
			<-backend.refreshChan
			backend.draw()
		}
	}()

	for {
		event := backend.tcellScreen.PollEvent()
		switch event := event.(type) {
		case *tcell.EventResize:
			backend.tcellScreen.Sync()
			backend.screen.eventChan <- Event(EventScreen{})
		case *tcell.EventKey:
			if backend.screen.UseKeys {
				backend.screen.eventChan <- Event(backend.tCellEventKeyToEventKey(event))
			}
		case *tcell.EventMouse:
			if backend.screen.UseMouse {
				x, y := event.Position()
				//button := event.Buttons()
				backend.screen.eventChan <- Event(EventMouse{X: x, Y: y})
			}
		}
		backend.draw()
	}
}

// Refresh causes the backend loop to redraw the screen.
func (backend *BackendTCell) Refresh() {
	backend.refreshChan <- struct{}{}
}

// Size returns the current backend window dimensions.
func (backend *BackendTCell) Size() (int, int) {
	return backend.tcellScreen.Size()
}

// SetSize sets the backend window to the provided width and height. Does nothing.
func (backend *BackendTCell) SetSize(w, h int) {
}

// Units returns the unit type the backend uses for Size().
func (backend *BackendTCell) Units() int {
	return UnitCells
}

// Scale returns the current backend window scaling. Does nothing.
func (backend *BackendTCell) Scale() float64 {
	return 0
}

// SetScale sets the backend window's scaling. Does nothing.
func (backend *BackendTCell) SetScale(scale float64) {
}

// draw is used for drawing the screen's cells to the tcell screen.
func (backend *BackendTCell) draw() {
	if backend.screen.Redraw {
		for y := 0; y < len(backend.screen.cells); y++ {
			for x := 0; x < len(backend.screen.cells[y]); x++ {
				if backend.screen.cells[y][x].Redraw {
					backend.tcellScreen.SetContent(x, y, backend.screen.cells[y][x].Rune, nil, StyleToTCellStyle(backend.screen.cells[y][x].Style))
					backend.screen.cells[y][x].Redraw = false
				}
			}
		}
		backend.screen.Redraw = false
	}
	backend.tcellScreen.Show()
}

func (backend *BackendTCell) tCellEventKeyToEventKey(tcellEvent *tcell.EventKey) (eventKey EventKey) {
	eventKey.Key = backend.tCellKeyToKey(tcellEvent.Key(), tcellEvent.Rune())
	eventKey.Rune = tcellEvent.Rune()
	modifiers := tcellEvent.Modifiers()
	eventKey.Ctrl = modifiers&tcell.ModCtrl != 0
	eventKey.Alt = modifiers&tcell.ModAlt != 0
	eventKey.Shift = modifiers&tcell.ModShift != 0
	// This is a little weird but I would like capitalized characters to report their shifted state...
	if _, ok := shiftMap[eventKey.Rune]; ok {
		eventKey.Shift = true
	}

	eventKey.Meta = modifiers&tcell.ModMeta != 0

	return
}

//
func (backend *BackendTCell) tCellKeyToKey(tcellKey tcell.Key, tcellRune rune) (key Key) {
	if tcellKey == tcell.KeyRune {
		var ok bool
		if key, ok = runeMap[tcellRune]; !ok {
			return KeyNull
		}
	} else {
		var ok bool
		if key, ok = keyMap[tcellKey]; !ok {
			return KeyNull
		}
	}
	return key
}

var keyMap = map[tcell.Key]Key{
	tcell.KeyF1:     KeyF1,
	tcell.KeyF2:     KeyF2,
	tcell.KeyF3:     KeyF3,
	tcell.KeyF4:     KeyF4,
	tcell.KeyF5:     KeyF5,
	tcell.KeyF6:     KeyF6,
	tcell.KeyF7:     KeyF7,
	tcell.KeyF8:     KeyF8,
	tcell.KeyF9:     KeyF9,
	tcell.KeyF10:    KeyF10,
	tcell.KeyF11:    KeyF11,
	tcell.KeyF12:    KeyF12,
	tcell.KeyLeft:   KeyLeft,
	tcell.KeyRight:  KeyRight,
	tcell.KeyUp:     KeyUp,
	tcell.KeyDown:   KeyDown,
	tcell.KeyEscape: KeyEscape,
}

var runeMap = map[rune]Key{
	' ':  KeySpace,
	'!':  KeyExclaim,
	'"':  KeyDoubleQuote,
	'#':  KeyHash,
	'$':  KeyDollar,
	'%':  KeyPercent,
	'&':  KeyAmpersand,
	'\'': KeyQuote,
	'(':  KeyLeftParenthesis,
	')':  KeyRightParenthesis,
	'*':  KeyAsterisk,
	'+':  KeyPlus,
	',':  KeyComma,
	'-':  KeyMinus,
	'.':  KeyPeriod,
	'/':  KeySlash,
	'0':  Key0,
	'1':  Key1,
	'2':  Key2,
	'3':  Key3,
	'4':  Key4,
	'5':  Key5,
	'6':  Key6,
	'7':  Key7,
	'8':  Key8,
	'9':  Key9,
	':':  KeyColon,
	';':  KeySemiColon,
	'<':  KeyLess,
	'=':  KeyEqual,
	'>':  KeyGreater,
	'?':  KeyQuestion,
	'`':  KeyBackquote,
	'@':  KeyAt,
	'a':  KeyA,
	'b':  KeyB,
	'c':  KeyC,
	'd':  KeyD,
	'e':  KeyE,
	'f':  KeyF,
	'g':  KeyG,
	'h':  KeyH,
	'i':  KeyI,
	'j':  KeyJ,
	'k':  KeyK,
	'l':  KeyL,
	'm':  KeyM,
	'n':  KeyN,
	'o':  KeyO,
	'p':  KeyP,
	'q':  KeyQ,
	'r':  KeyR,
	's':  KeyS,
	't':  KeyT,
	'u':  KeyU,
	'v':  KeyV,
	'w':  KeyW,
	'x':  KeyX,
	'y':  KeyY,
	'z':  KeyZ,
	'A':  KeyA,
	'B':  KeyB,
	'C':  KeyC,
	'D':  KeyD,
	'E':  KeyE,
	'F':  KeyF,
	'G':  KeyG,
	'H':  KeyH,
	'I':  KeyI,
	'J':  KeyJ,
	'K':  KeyK,
	'L':  KeyL,
	'M':  KeyM,
	'N':  KeyN,
	'O':  KeyO,
	'P':  KeyP,
	'Q':  KeyQ,
	'R':  KeyR,
	'S':  KeyS,
	'T':  KeyT,
	'U':  KeyU,
	'V':  KeyV,
	'W':  KeyW,
	'X':  KeyX,
	'Y':  KeyY,
	'Z':  KeyZ,
}

var shiftMap = map[rune]bool{
	'A': true,
	'B': true,
	'C': true,
	'D': true,
	'E': true,
	'F': true,
	'G': true,
	'H': true,
	'I': true,
	'J': true,
	'K': true,
	'L': true,
	'M': true,
	'N': true,
	'O': true,
	'P': true,
	'Q': true,
	'R': true,
	'S': true,
	'T': true,
	'U': true,
	'V': true,
	'W': true,
	'X': true,
	'Y': true,
	'Z': true,
}
