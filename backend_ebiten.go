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
	"path"
	"strings"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/text"
	"github.com/kettek/goro/glyphs"
	"github.com/kettek/goro/resources"
)

// BackendEbiten is our Ebiten backend.
type BackendEbiten struct {
	screen        Screen
	imageBuffer   *ebiten.Image
	op            *ebiten.DrawImageOptions
	width, height int
	hasStarted    bool
	glyphs        []glyphs.Glyphs
	emptyCell     *ebiten.Image

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

	backend.imageBuffer, _ = ebiten.NewImage(320, 240, ebiten.FilterDefault)
	backend.op = &ebiten.DrawImageOptions{}

	backend.glyphs = make([]glyphs.Glyphs, 10)
	backend.emptyCell, _ = ebiten.NewImage(16, 16, ebiten.FilterDefault)

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
	err = ebiten.Run(func(screenBuffer *ebiten.Image) (err error) {
		if !backend.hasStarted {
			// Load our built-in glyphs
			backend.SetGlyphsFromTTFBytes(0, resources.GoroTTF, 16)

			go func() {
				cb(&backend.screen)
			}()
			backend.hasStarted = true
		}

		keyEvents := make([]EventKey, 0)
		// ... Ew.
		for k := ebiten.Key(0); k <= ebiten.KeyMax; k++ {
			if ebiten.IsKeyPressed(k) {
				if !backend.pressedKeys[k] {
					keyEvents = append(keyEvents, backend.ebitenKeyToEventKey(k))
				}
				backend.pressedKeys[k] = true
			} else {
				backend.pressedKeys[k] = false
			}
		}
		// FIXME: This isn't exactly right for non-US keyboards...
		inputRunes := ebiten.InputChars()
		for i, k := range keyEvents {
			for i2, r := range inputRunes {
				if k2, ok := RuneToKeyMap[r]; ok {
					if k2 == k.Key {
						keyEvents[i].Rune = r
						inputRunes = append(inputRunes[:i2], inputRunes[i2+1:]...)
						break
					}
				}
			}
		}
		// Convert remaining inputRunes to KeyEvents
		for _, r := range inputRunes {
			keyEvents = append(keyEvents, EventKey{
				Key:   KeyNull,
				Rune:  r,
				Shift: backend.pressedKeys[ebiten.KeyShift],
				Ctrl:  backend.pressedKeys[ebiten.KeyControl],
				Alt:   backend.pressedKeys[ebiten.KeyAlt],
			})
		}
		// Send our KeyEvents
		if backend.screen.UseKeys {
			for _, k := range keyEvents {
				backend.screen.eventChan <- k
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
		if !ebiten.IsDrawingSkipped() {
			if backend.screen.Redraw {
				backend.drawCellBackgrounds(backend.imageBuffer)
				backend.drawCellForegrounds(backend.imageBuffer)
				backend.screen.Redraw = false
			}

			backend.op.GeoM.Reset()
			screenBuffer.DrawImage(backend.imageBuffer, backend.op)
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
	backend.imageBuffer, _ = ebiten.NewImage(w, h, ebiten.FilterDefault)
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

// SetGlyphs sets the glyphs to be used for rendering.
func (backend *BackendEbiten) SetGlyphs(id glyphs.ID, filePath string, size float64) error {
	ext := strings.ToLower(path.Ext(filePath))
	switch ext {
	case ".ttf":
		{
			ttfGlyphs, err := glyphs.LoadTruetype(filePath)
			if err != nil {
				return err
			}
			ttfGlyphs.SetSize(size)
			backend.glyphs[id] = ttfGlyphs
		}
	default:
		return nil
	}
	backend.syncGlyphs(id)
	return nil
}

// SetGlyphsFromBytes sets the glyphs to be used for rendering from the provided TTF bytes.
func (backend *BackendEbiten) SetGlyphsFromTTFBytes(id glyphs.ID, bytes []byte, size float64) error {
	ttfGlyphs, err := glyphs.LoadTruetypeFromBytes(bytes)
	if err != nil {
		return err
	}
	ttfGlyphs.SetSize(size)
	backend.glyphs[id] = ttfGlyphs
	backend.syncGlyphs(id)
	return nil
}

// syncGlyphs synchronizes the screen's size and backend size, along with associated cached variables, to use the updated glyphs.
func (backend *BackendEbiten) syncGlyphs(id glyphs.ID) {
	c, r := backend.screen.Size()
	newWidth := c * backend.glyphs[id].Width()
	newHeight := r * backend.glyphs[id].Height()

	backend.SetSize(newWidth, newHeight)

	// Is this needed?
	backend.emptyCell, _ = ebiten.NewImage(backend.glyphs[id].Width(), backend.glyphs[id].Height(), ebiten.FilterDefault)

	backend.screen.ForceRedraw()
	backend.Refresh()
}

// SyncSize is the external call to synchronize the screen's size.
func (backend *BackendEbiten) SyncSize() {
	c, r := backend.screen.Size()
	newWidth := c * backend.glyphs[0].Width()
	newHeight := r * backend.glyphs[0].Height()

	backend.SetSize(newWidth, newHeight)

	backend.Refresh()
}

// drawCellForegrounds draws the colored glyphs for the cell at x and y.
func (backend *BackendEbiten) drawCellForegrounds(target *ebiten.Image) {
	backend.screen.cellsMutex.Lock()
	for y := 0; y < len(backend.screen.cells); y++ {
		for x := 0; x < len(backend.screen.cells[y]); x++ {
			if !backend.screen.cells[y][x].Redraw {
				continue
			}
			fg := backend.screen.cells[y][x].Style.Foreground
			if fg == ColorNone {
				fg = backend.screen.Foreground
			}
			// Draw our rune
			if backend.screen.cells[y][x].Rune != rune(0) {
				glyphSet := backend.glyphs[backend.screen.cells[y][x].Glyphs]
				switch glyphSet := glyphSet.(type) {
				case *glyphs.Truetype:
					bounds, _, _ := glyphSet.Normal.GlyphBounds(backend.screen.cells[y][x].Rune)
					text.Draw(
						target,
						string(backend.screen.cells[y][x].Rune),
						glyphSet.Normal,
						x*glyphSet.Width()+(glyphSet.Width()/2-bounds.Max.X.Round()/2),
						y*glyphSet.Height()+glyphSet.Ascent(),
						fg,
					)
				}
			}
			backend.screen.cells[y][x].Redraw = false
		}
	}
	backend.screen.cellsMutex.Unlock()
}

// drawCellBackgrounds draws the background at the cell at x and y
func (backend *BackendEbiten) drawCellBackgrounds(target *ebiten.Image) {
	backend.screen.cellsMutex.Lock()
	for y := 0; y < len(backend.screen.cells); y++ {
		for x := 0; x < len(backend.screen.cells[y]); x++ {
			if !backend.screen.cells[y][x].Redraw {
				continue
			}
			bg := backend.screen.cells[y][x].Style.Background
			if bg == ColorNone {
				bg = backend.screen.Background
			}
			id := backend.screen.cells[y][x].Glyphs
			backend.emptyCell.Fill(bg)
			backend.op.GeoM.Reset()
			backend.op.GeoM.Translate(float64(x*backend.glyphs[id].Width()), float64(y*backend.glyphs[id].Height()))
			target.DrawImage(backend.emptyCell, backend.op)
		}
	}
	backend.screen.cellsMutex.Unlock()
}

func (backend *BackendEbiten) DrawRect(image *ebiten.Image, x0, y0, x1, y1 float32, c Color) {
	r := float32(c.R) / 0xff
	g := float32(c.G) / 0xff
	b := float32(c.B) / 0xff
	a := float32(c.A) / 0xff

	vertices := []ebiten.Vertex{
		{
			DstX: x0, DstY: y0,
			SrcX: 1, SrcY: 1,
			ColorR: r, ColorG: g, ColorB: b, ColorA: a,
		},
		{
			DstX: x1, DstY: y0,
			SrcX: 1, SrcY: 1,
			ColorR: r, ColorG: g, ColorB: b, ColorA: a,
		},
		{
			DstX: x0, DstY: y1,
			SrcX: 1, SrcY: 1,
			ColorR: r, ColorG: g, ColorB: b, ColorA: a,
		},
		{
			DstX: x1, DstY: y1,
			SrcX: 1, SrcY: 1,
			ColorR: r, ColorG: g, ColorB: b, ColorA: a,
		},
	}
	indices := []uint16{0, 1, 2, 1, 2, 3}

	image.DrawTriangles(vertices, indices, backend.emptyCell, nil)
}

func (backend *BackendEbiten) ebitenKeyToEventKey(k ebiten.Key) (eventKey EventKey) {
	var key Key
	var ok bool
	if key, ok = ebitenKeyMap[k]; !ok {
		key = KeyNull
	}

	eventKey.Key = key

	if backend.pressedKeys[ebiten.KeyShift] {
		eventKey.Shift = true
	}
	if backend.pressedKeys[ebiten.KeyControl] {
		eventKey.Ctrl = true
	}
	if backend.pressedKeys[ebiten.KeyAlt] {
		eventKey.Alt = true
	}
	// TODO: Meta?

	return
}

var ebitenKeyMap = map[ebiten.Key]Key{
	ebiten.Key0:            Key0,
	ebiten.Key1:            Key1,
	ebiten.Key2:            Key2,
	ebiten.Key3:            Key3,
	ebiten.Key4:            Key4,
	ebiten.Key5:            Key5,
	ebiten.Key6:            Key6,
	ebiten.Key7:            Key7,
	ebiten.Key8:            Key8,
	ebiten.Key9:            Key9,
	ebiten.KeyA:            KeyA,
	ebiten.KeyB:            KeyB,
	ebiten.KeyC:            KeyC,
	ebiten.KeyD:            KeyD,
	ebiten.KeyE:            KeyE,
	ebiten.KeyF:            KeyF,
	ebiten.KeyG:            KeyG,
	ebiten.KeyH:            KeyH,
	ebiten.KeyI:            KeyI,
	ebiten.KeyJ:            KeyJ,
	ebiten.KeyK:            KeyK,
	ebiten.KeyL:            KeyL,
	ebiten.KeyM:            KeyM,
	ebiten.KeyN:            KeyN,
	ebiten.KeyO:            KeyO,
	ebiten.KeyP:            KeyP,
	ebiten.KeyQ:            KeyQ,
	ebiten.KeyR:            KeyR,
	ebiten.KeyS:            KeyS,
	ebiten.KeyT:            KeyT,
	ebiten.KeyU:            KeyU,
	ebiten.KeyV:            KeyV,
	ebiten.KeyW:            KeyW,
	ebiten.KeyX:            KeyX,
	ebiten.KeyY:            KeyY,
	ebiten.KeyZ:            KeyZ,
	ebiten.KeyAlt:          KeyAlt,
	ebiten.KeyApostrophe:   KeyApostrophe,
	ebiten.KeyBackslash:    KeyBackslash,
	ebiten.KeyBackspace:    KeyBackspace,
	ebiten.KeyCapsLock:     KeyCapsLock,
	ebiten.KeyComma:        KeyComma,
	ebiten.KeyControl:      KeyControl,
	ebiten.KeyDelete:       KeyDelete,
	ebiten.KeyEnd:          KeyEnd,
	ebiten.KeyEnter:        KeyEnter,
	ebiten.KeyEqual:        KeyEqual,
	ebiten.KeyEscape:       KeyEscape,
	ebiten.KeyF1:           KeyF1,
	ebiten.KeyF2:           KeyF2,
	ebiten.KeyF3:           KeyF3,
	ebiten.KeyF4:           KeyF4,
	ebiten.KeyF5:           KeyF5,
	ebiten.KeyF6:           KeyF6,
	ebiten.KeyF7:           KeyF7,
	ebiten.KeyF8:           KeyF8,
	ebiten.KeyF9:           KeyF9,
	ebiten.KeyF10:          KeyF10,
	ebiten.KeyF11:          KeyF11,
	ebiten.KeyF12:          KeyF12,
	ebiten.KeyGraveAccent:  KeyGraveAccent,
	ebiten.KeyHome:         KeyHome,
	ebiten.KeyInsert:       KeyInsert,
	ebiten.KeyKP0:          KeyKP0,
	ebiten.KeyKP1:          KeyKP1,
	ebiten.KeyKP2:          KeyKP2,
	ebiten.KeyKP3:          KeyKP3,
	ebiten.KeyKP4:          KeyKP4,
	ebiten.KeyKP5:          KeyKP5,
	ebiten.KeyKP6:          KeyKP6,
	ebiten.KeyKP7:          KeyKP7,
	ebiten.KeyKP8:          KeyKP8,
	ebiten.KeyKP9:          KeyKP9,
	ebiten.KeyKPAdd:        KeyKPAdd,
	ebiten.KeyKPDecimal:    KeyKPDecimal,
	ebiten.KeyKPDivide:     KeyKPDivide,
	ebiten.KeyKPEnter:      KeyKPEnter,
	ebiten.KeyKPEqual:      KeyKPEqual,
	ebiten.KeyKPMultiply:   KeyKPMultiply,
	ebiten.KeyKPSubtract:   KeyKPSubtract,
	ebiten.KeyLeftBracket:  KeyLeftBracket,
	ebiten.KeyMenu:         KeyMenu,
	ebiten.KeyMinus:        KeyMinus,
	ebiten.KeyNumLock:      KeyNumLock,
	ebiten.KeyPageDown:     KeyPageDown,
	ebiten.KeyPageUp:       KeyPageUp,
	ebiten.KeyPause:        KeyPause,
	ebiten.KeyPeriod:       KeyPeriod,
	ebiten.KeyPrintScreen:  KeyPrintScreen,
	ebiten.KeyRightBracket: KeyRightBracket,
	ebiten.KeyScrollLock:   KeyScrollLock,
	ebiten.KeySemicolon:    KeySemicolon,
	ebiten.KeyShift:        KeyShift,
	ebiten.KeySlash:        KeySlash,
	ebiten.KeySpace:        KeySpace,
	ebiten.KeyTab:          KeyTab,
	ebiten.KeyLeft:         KeyLeft,
	ebiten.KeyRight:        KeyRight,
	ebiten.KeyUp:           KeyUp,
	ebiten.KeyDown:         KeyDown,
}
