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

// Key is our type representing keyboard codes.
type Key uint8

// These are our constant key variables.
const (
	KeyNull Key = iota
	_
	_
	_
	_
	_
	_
	_
	KeyBackspace
	KeyTab
	_
	_
	_
	KeyReturn
	_
	_
	_
	_
	_
	_
	_
	_
	_
	_
	_
	_
	_
	KeyEscape
	_
	_
	_
	_
	KeySpace
	KeyExclaim
	KeyDoubleQuote
	KeyHash
	KeyDollar
	KeyPercent
	KeyAmpersand
	KeyQuote
	KeyLeftParenthesis
	KeyRightParenthesis
	KeyAsterisk
	KeyPlus
	KeyComma
	KeyMinus
	KeyPeriod
	KeySlash
	Key0
	Key1
	Key2
	Key3
	Key4
	Key5
	Key6
	Key7
	Key8
	Key9
	KeyColon
	KeySemiColon
	KeyLess
	KeyEqual
	KeyGreater
	KeyQuestion
	KeyAt
	KeyLeftBracket = iota + 26
	KeyBackslash
	KeyRightBracket
	KeyCaret
	KeyUnderscore
	KeyBackquote
	KeyA
	KeyB
	KeyC
	KeyD
	KeyE
	KeyF
	KeyG
	KeyH
	KeyI
	KeyJ
	KeyK
	KeyL
	KeyM
	KeyN
	KeyO
	KeyP
	KeyQ
	KeyR
	KeyS
	KeyT
	KeyU
	KeyV
	KeyW
	KeyX
	KeyY
	KeyZ
	KeyLeftCurlyBracket
	KeyPipe
	KeyRightCurlyBracket
	KeyTilde
	KeyDelete
	KeyAlt
	KeyControl
	KeyShift
	KeyCapsLock
	KeyEnd
	KeyMenu
	KeyPrintScreen
	KeyScrollLock
	KeyEnter
	KeyInsert
	KeyF1
	KeyF2
	KeyF3
	KeyF4
	KeyF5
	KeyF6
	KeyF7
	KeyF8
	KeyF9
	KeyF10
	KeyF11
	KeyF12
	KeyKP0
	KeyKP1
	KeyKP2
	KeyKP3
	KeyKP4
	KeyKP5
	KeyKP6
	KeyKP7
	KeyKP8
	KeyKP9
	KeyKPAdd
	KeyKPPeriod
	KeyKPDivide
	KeyKPEnter
	KeyKPEquals
	KeyKPMultiply
	KeyKPSubtract
	KeyLeft
	KeyUp
	KeyRight
	KeyDown
)

// Mod represents a modifier key state.
type Mod uint8

// These are our const modifier variables.
const (
	ModShift Mod = 1 << iota
	ModCtrl
	ModAlt
	ModMeta
	ModNone Mod = 0
)
