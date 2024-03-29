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

import (
	"math"
)

//
type NodeAStar struct {
	parent              *NodeAStar
	x, y                int // for deconstructing the path
	gCost, hCost, fCost float64
	mCost               uint32
}

// PathAStar represents a pathing structure for the A* algorithm.
type PathAStar struct {
	width, height  int
	nodes          [][]*NodeAStar
	heuristicsFunc func(x0, y0 int, x1, y1 int) float64
	diagonals      bool
}

func NewPathAStarFromMap(pathMap PathMap) Path {
	path := &PathAStar{}
	path.Resize(pathMap.Width(), pathMap.Height())

	for y, n := range path.nodes {
		for x, _ := range n {
			path.nodes[y][x] = &NodeAStar{
				y:     y,
				x:     x,
				fCost: math.MaxFloat64,
				gCost: math.MaxFloat64,
				hCost: math.MaxFloat64,
				mCost: pathMap.CostAt(x, y),
			}
		}
	}

	// calculate from each tile's movement cost.
	// calculate Blocked, etc. ?
	return path
}

func NewPathAStarFromFunc(width, height int, calcFunc func(int, int) uint32) Path {
	path := &PathAStar{}
	path.Resize(width, height)

	for y, n := range path.nodes {
		for x, _ := range n {
			path.nodes[y][x] = &NodeAStar{
				y:     y,
				x:     x,
				fCost: math.MaxFloat64,
				gCost: math.MaxFloat64,
				hCost: math.MaxFloat64,
				mCost: calcFunc(x, y),
			}
		}
	}
	return path
}

// Resize resizes the given MapBase to the provided size.
func (p *PathAStar) Resize(width, height int) {
	p.width = width
	p.height = height

	currHeight := len(p.nodes)
	// Grow or shrink our height.
	if currHeight < p.height {
		p.nodes = append(p.nodes, make([][]*NodeAStar, p.height-currHeight)...)
	} else if currHeight > p.height {
		p.nodes = p.nodes[:p.height]
	}
	// Iterate through our height to grow or shrink their width.
	for y := range p.nodes {
		currWidth := len(p.nodes[y])
		if currWidth < p.width {
			p.nodes[y] = append(p.nodes[y], make([]*NodeAStar, p.width-currWidth)...)
		} else if currWidth > p.width {
			p.nodes[y] = p.nodes[y][:p.width]
		}
	}
}

func (p *PathAStar) Compute(oX, oY int, tX, tY int) (steps []Step) {
	// Sanity checks.
	if tX < 0 || tX >= p.width || tY < 0 || tY >= p.height {
		return
	}
	if oX == tX && oY == tY {
		return
	}
	// Set our first node's costs.
	p.nodes[oY][oX].gCost = 0
	p.nodes[oY][oX].fCost = p.calculateH(oX, oY, tX, tY)

	// Create our open nodes slice.
	openNodes := make([]*NodeAStar, 0)
	openNodes = append(openNodes, p.nodes[oY][oX])

	for len(openNodes) > 0 {
		index := 0
		var current *NodeAStar = openNodes[0]
		// Get node with lowest fCost
		for i, n := range openNodes {
			if n.fCost < current.fCost {
				current = n
				index = i
			}
		}

		// If it is our destination then we've found a path.
		if current.y == tY && current.x == tX {
			steps = p.tracePath(tX, tY)
			return
		}

		// Remove it.
		openNodes = append(openNodes[:index], openNodes[index+1:]...)

		var steps [][2]int
		if p.diagonals {
			steps = [][2]int{
				{-1, -1}, //tl
				{0, -1},  // t
				{1, -1},  // tr
				{-1, 0},  // l
				//{0, 0},   // c
				{1, 0},  // r
				{-1, 1}, // bl
				{0, 1},  // b
				{1, 1},  // br
			}
		} else {
			steps = [][2]int{
				{0, -1}, // t
				{0, 1},  // b
				//{0, 0},  // c
				{-1, 0}, // l
				{1, 0},  // r
			}
		}
		// Iterate through our eight neighboring nodes.
		for _, ij := range steps {
			i := ij[0]
			j := ij[1]
			// Sanity checks
			if current.y+i < 0 || current.y+i >= p.height || current.x+j < 0 || current.x+j >= p.width {
				continue
			}
			neighbor := p.nodes[current.y+i][current.x+j]
			// Skip neighbor if it has maximum cost aka blocking
			if neighbor.mCost == MaximumCost {
				continue
			}
			g := current.gCost + 1 + float64(neighbor.mCost)
			// Add extra diagonal cost.
			if p.diagonals && (math.Abs(float64(i))+math.Abs(float64(j))) == 2 {
				g += .414
			}

			// The path is better. Record it.
			if g < neighbor.gCost {
				neighbor.parent = current
				neighbor.gCost = g
				neighbor.fCost = g + p.calculateH(current.x+i, current.y+j, tX, tY)
				exists := false
				for _, node := range openNodes {
					if node == neighbor {
						exists = true
					}
				}
				if !exists {
					openNodes = append(openNodes, neighbor)
				}
			}
		}
		//
	}
	return
}

func (p *PathAStar) calculateH(x, y int, tX, tY int) float64 {
	if p.heuristicsFunc != nil {
		return p.heuristicsFunc(x, y, tX, tY)
	}
	// Euclidean Distance
	a := math.Pow(float64(y-tY), 2)
	b := math.Pow(float64(x-tX), 2)
	return math.Sqrt(a + b)
}

func (p *PathAStar) tracePath(tX, tY int) (steps []Step) {
	for node := p.nodes[tY][tX]; node.parent != nil; node = node.parent {
		steps = append(steps, Step{x: node.x, y: node.y})
	}

	// Reverse our steps.
	for i, j := 0, len(steps)-1; i < j; i, j = i+1, j-1 {
		steps[i], steps[j] = steps[j], steps[i]
	}

	return
}

func (p *PathAStar) SetHeuristicsFunc(f func(x0, y0 int, x1, y1 int) float64) {
	p.heuristicsFunc = f
}

// AllowDiagonals sets if diagonal movement is allowed.
func (p *PathAStar) AllowDiagonals(v bool) {
	p.diagonals = v
}
