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
	Reset()
	Recompute(cX, cY int, radius int, light Light)
	Compute(cX, cY int, radius int, light Light)
	SetCell(x, y int, fovCell Cell) error
	BlocksMovement(x, y int) bool
	SetBlocksMovement(x, y int, blocks bool) error
	BlocksLight(x, y int) bool
	SetBlocksLight(x, y int, blocks bool) error
	Visible(x, y int) bool
	SetVisible(x, y int, visible bool) error
	Lighting(x, y int) Light
	SetLighting(x, y int, light Light) error
	CheckBounds(x, y int) error
	ToString(showVisible, showBlocksLight, showBlocksMovement, showLight bool) string
}

// Algorithm represents a field-of-view algorithm
type Algorithm uint8

// Our default algorithms
const (
	AlgorithmNone Algorithm = iota
	AlgorithmBBQ
)

// NewMap returns a new Map at the given dimensions using the provided Algorithm.
func NewMap(width, height int, algo Algorithm) Map {
	switch algo {
	default:
		return nil
	case AlgorithmNone:
		return nil
	case AlgorithmBBQ:
		return NewMapBBQ(width, height)
	}
}
