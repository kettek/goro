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

import (
  "math"
	"errors"
	"strconv"
)

// MapBase is the base implementation of our Map interface. It is used as an embedded structure within most Map implementations.
type MapBase struct {
	width, height int
	cells         [][]Cell
}

// Resize resizes the given MapBase to the provided size.
func (fovMap *MapBase) Resize(width, height int) {
	fovMap.width = width
	fovMap.height = height
	currHeight := len(fovMap.cells)
	// Grow or shrink our height.
	if currHeight < fovMap.height {
		fovMap.cells = append(fovMap.cells, make([][]Cell, fovMap.height-currHeight)...)
	} else if currHeight > fovMap.height {
		fovMap.cells = fovMap.cells[:fovMap.height]
	}
	// Iterate through our height to grow or shrink their width.
	for y := range fovMap.cells {
		currWidth := len(fovMap.cells[y])
		if currWidth < fovMap.width {
			fovMap.cells[y] = append(fovMap.cells[y], make([]Cell, fovMap.width-currWidth)...)
		} else if currWidth > fovMap.width {
			fovMap.cells[y] = fovMap.cells[y][:fovMap.width]
		}
	}
}

// Clear clears the map using the provided Cell for clearing.
func (fovMap *MapBase) Clear(fovCell Cell) {
	for y := range fovMap.cells {
		for x := range fovMap.cells[y] {
			fovMap.cells[y][x] = fovCell
		}
	}
}

// SetCell sets the cell at the given x and y to fovCell.
func (fovMap *MapBase) SetCell(x, y int, fovCell Cell) error {
	if err := fovMap.CheckBounds(x, y); err != nil {
		return err
	}
	fovMap.cells[y][x] = fovCell
	return nil
}

// BlocksMovement checks if the cell at x and y blocks movement. Returns true if x or y is out of bounds.
func (fovMap *MapBase) BlocksMovement(x, y int) bool {
	if err := fovMap.CheckBounds(x, y); err != nil {
		return true
	}
	return fovMap.cells[y][x].BlocksMovement
}

// SetBlocksMovement sets the blocksMovement property of the cell at x and y to blocks.
func (fovMap *MapBase) SetBlocksMovement(x, y int, blocks bool) error {
	if err := fovMap.CheckBounds(x, y); err != nil {
		return err
	}
	fovMap.cells[y][x].BlocksMovement = blocks
	return nil
}

// BlocksLight checks if the cell at x and y blocks light. Returns true if x or y is out of bounds.
func (fovMap *MapBase) BlocksLight(x, y int) bool {
	if err := fovMap.CheckBounds(x, y); err != nil {
		return true
	}
	return fovMap.cells[y][x].BlocksLight
}

// SetBlocksLight sets the blocksLight property of the cell at x and y to blocks.
func (fovMap *MapBase) SetBlocksLight(x, y int, blocks bool) error {
	if err := fovMap.CheckBounds(x, y); err != nil {
		return err
	}
	fovMap.cells[y][x].BlocksLight = blocks
	return nil
}

// Visible returns whether a given cell is considered as within the FoV. Returns false if x or y is out of bounds.
func (fovMap *MapBase) Visible(x, y int) bool {
	if err := fovMap.CheckBounds(x, y); err != nil {
		return false
	}
	return fovMap.cells[y][x].Visible
}

// SetVisible sets whether a given cell is considered as within the FoV.
func (fovMap *MapBase) SetVisible(x, y int, visible bool) error {
	if err := fovMap.CheckBounds(x, y); err != nil {
		return err
	}
	fovMap.cells[y][x].Visible = visible
	return nil
}

// Lighting returns the lighting value for a given cell. Returns 0 if x or y is out of bounds.
func (fovMap *MapBase) Lighting(x, y int) Light {
	if err := fovMap.CheckBounds(x, y); err != nil {
		return Light{}
	}
	return fovMap.cells[y][x].Lighting
}

// SetLighting sets a given cell's lighting value to the one provided.
func (fovMap *MapBase) SetLighting(x, y int, light Light) error {
	if err := fovMap.CheckBounds(x, y); err != nil {
		return err
	}
	fovMap.cells[y][x].Lighting = light
	return nil
}

// CheckBounds returns an error if the provided x and y is out of bounds.
func (fovMap *MapBase) CheckBounds(x, y int) error {
	if y < 0 || y >= len(fovMap.cells) {
		return errors.New("y out of range")
	}
	if x < 0 || x >= len(fovMap.cells[y]) {
		return errors.New("x out of range")
	}
	return nil
}

// Reset resets the visibility state of all cells.
func (fovMap *MapBase) Reset() {
	for y := range fovMap.cells {
		for x := range fovMap.cells[y] {
			fovMap.cells[y][x].Visible = false
		}
	}
}

// Width returns the width of the map.
func (fovMap *MapBase) Width() int {
  return fovMap.width
}

// Height returns the height of the map.
func (fovMap *MapBase) Height() int {
  return fovMap.height
}

// Height returns the movement cost of a tile in a map. This allows fov.Maps to be used for pathing.
func (fovMap *MapBase) CostAt(x, y int) uint32 {
  if err := fovMap.CheckBounds(x, y); err != nil {
    return math.MaxUint32
  }
  cost := 0
	if fovMap.cells[y][x].BlocksMovement {
    cost = math.MaxUint32
  } else {
    if fovMap.cells[y][x].BlocksLight {
      cost++
    }
    if !fovMap.cells[y][x].Visible {
      cost++
    }
  }
  return cost
}

// ToString returns a stringified view of the map.
func (fovMap *MapBase) ToString(showVisible, showBlocksLight, showBlocksMovement, showLight bool) string {
	var str string
	for y := range fovMap.cells {
		for x := range fovMap.cells[y] {
			if showVisible {
				if fovMap.cells[y][x].Visible {
					str += "v"
				} else {
					str += " "
				}
			}
			if showBlocksLight {
				if fovMap.cells[y][x].BlocksLight {
					str += "B"
				} else {
					str += " "
				}
			}
			if showBlocksMovement {
				if fovMap.cells[y][x].BlocksMovement {
					str += "b"
				} else {
					str += " "
				}
			}
			if showLight {
				str += strconv.Itoa(int(fovMap.cells[y][x].Lighting.Lumens))
			}
		}
		str += "\n"
	}
	return str
}
