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

package fov

// Map is our interface for field of view Maps.
type Map interface {
	Resize(width, height int)
	Clear(fovCell Cell)
	Compute(cX, cY int, radius int, light int8)
	SetCell(x, y int, fovCell Cell) error
	BlocksMovement(x, y int) bool
	SetBlocksMovement(x, y int, blocks bool) error
	BlocksLight(x, y int) bool
	SetBlocksLight(x, y int, blocks bool) error
	Seen(x, y int) bool
	SetSeen(x, y int, seen bool) error
	Lighting(x, y int) int8
	SetLighting(x, y int, light int8) error
	CheckBounds(x, y int) error
}

type Algorithm uint8

const (
	AlgorithmNone Algorithm = iota
	AlgorithmBBQ
)

func NewMap(width, height int, algo Algorithm) Map {
	switch algo {
	default:
		return nil
	case AlgorithmNone:
		return nil
	case AlgorithmBBQ:
		return NewMapBBQ(width, height)
	}
	return nil
}
