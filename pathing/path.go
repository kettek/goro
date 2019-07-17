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
package pathing

// Path is our interface for paths.
type Path interface {
  Compute(oX, oY int, tX, tY int) []Step
  SetHeuristicsFunc(func(x0, y0 int, x1, y1 int) float64)
}

// Algorithm represents a pathing algorithm
type Algorithm uint8

// Our default algorithms
const (
	AlgorithmNone Algorithm = iota
  AlgorithmAStar
)

// NewPathFromMap returns a new pathing map from the given pathMap interface.
func NewPathFromMap(pathMap PathMap, algo Algorithm) Path {
	switch algo {
	default:
		return nil
	case AlgorithmNone:
		return nil
	case AlgorithmAStar:
		return NewPathAStarFromMap(pathMap)
	}
}

// NewPathFromFunc returns a new pathing map with movement costs calculated from calcFunc.
func NewPathFromFunc(width, height int, calcFunc func(x int, y int) (cost int), algo Algorithm) Path {
	switch algo {
	default:
		return nil
	case AlgorithmNone:
		return nil
	case AlgorithmAStar:
		return NewPathAStarFromFunc(width, height, calcFunc)
	}
}
