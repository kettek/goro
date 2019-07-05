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
)

// Map represents a 2D structure for calculating a field of view.
type MapBBQ struct {
	MapBase
}

// NewMapBBQ returns an field of view Map sized to the given width and height using the BBQ algorithm.
func NewMapBBQ(width, height int) (fovMap *MapBBQ) {
	fovMap = &MapBBQ{}
	fovMap.Resize(width, height)
	return fovMap
}

// Compute calculates the FOV for our BBQ.
func (fovMap *MapBBQ) Compute(cX, cY int, radius int, light Light) {
	maxRadius := math.Sqrt(float64(radius*radius + radius*radius))
	for i := -radius; i <= radius; i++ {
		for j := -radius; j <= radius; j++ {
			if i*i+j*j < radius*radius {
				fovMap.computeLOS(cX, cY, cX+i, cY+j, maxRadius, light)
			}
		}
	}
}

// computeLOS checks the line of sight between x0,y0 to x1,y1, setting the cell at x1,y1 to
func (fovMap *MapBBQ) computeLOS(x0, y0, x1, y1 int, maxRadius float64, light Light) {
	var quadrantX, quadrantY, nextX, nextY, destX, destY int
	var distance float64

	destX = x1 - x0
	destY = y1 - y0

	if x0 < x1 {
		quadrantX = 1
	} else {
		quadrantX = -1
	}

	if y0 < y1 {
		quadrantY = 1
	} else {
		quadrantY = -1
	}

	nextX = x0
	nextY = y0

	distance = math.Sqrt(float64(destX*destX + destY*destY))

	for nextX != x1 || nextY != y1 {
		if fovMap.CheckBounds(nextX, nextY) != nil {
			return
		}
		if fovMap.BlocksLight(nextX, nextY) {
			return
		}

		if (math.Abs(float64(destY*(nextX-x0+quadrantX)-destX*(nextY-y0))) / distance) < 0.5 {
			nextX += quadrantX
		} else if (math.Abs(float64(destY*(nextX-x0)-destX*(nextY-y0+quadrantY))) / distance) < 0.5 {
			nextY += quadrantY
		} else {
			nextX += quadrantX
			nextY += quadrantY
		}
	}
	if fovMap.CheckBounds(x1, y1) != nil {
		return
	}
	fovMap.cells[y1][x1].Visible = true
	// Do light calculations (?)
	// It might be best if we make a separate light-only FoV calculation that aggregates multiple FoVs (from each light) and then concatenates their values together in a final FoV. But, for now, this works for a player-only system.
	lumens := light.Lumens * int16(distance/maxRadius)
	if lumens != 0 {
		if fovMap.cells[y1][x1].Lighting.Lumens < lumens {
			fovMap.cells[y1][x1].Lighting.Lumens += lumens
		}
	}

}
