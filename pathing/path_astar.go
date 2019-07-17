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
  "fmt"
	"math"
  "errors"
)

//
type NodeAStar struct {
  x,y int
  parentX, parentY int
  mCost int
  gCost, hCost, fCost float64
}

// PathAStar represents a pathing structure for the A* algorithm.
type PathAStar struct {
  width, height int
  foundPath bool
  pathMap PathMap
  nodes [][]NodeAStar
}

func NewPathAStarFromMap(pathMap PathMap) Path {
  path := &PathAStar{
    pathMap: pathMap,
  }
  path.Resize(pathMap.Width(), pathMap.Height())

  for y, n := range path.nodes {
    for x, _ := range n {
      path.nodes[y][x] = NodeAStar{
        fCost: math.MaxFloat64,
        gCost: math.MaxFloat64,
        hCost: math.MaxFloat64,
        mCost: pathMap.CostAt(x, y),
        parentX: -1,
        parentY: -1,
        x: x,
        y: y,
      }
    }
  }

  // calculate from each tile's movement cost.
  // calculate Blocked, etc. ?
  return path
}

// Resize resizes the given MapBase to the provided size.
func (p *PathAStar) Resize(width, height int) {
  p.width = width
  p.height = height

	currHeight := len(p.nodes)
	// Grow or shrink our height.
	if currHeight < p.height {
		p.nodes = append(p.nodes, make([][]NodeAStar, p.height-currHeight)...)
	} else if currHeight > p.height {
		p.nodes = p.nodes[:p.height]
	}
	// Iterate through our height to grow or shrink their width.
	for y := range p.nodes {
		currWidth := len(p.nodes[y])
		if currWidth < p.width {
			p.nodes[y] = append(p.nodes[y], make([]NodeAStar, p.width-currWidth)...)
		} else if currWidth > p.width {
			p.nodes[y] = p.nodes[y][:p.width]
		}
	}
}

func (p *PathAStar) Compute(oX, oY int, tX, tY int) error {
  closedList := make([][]bool, p.height)
  for y := range closedList {
    closedList[y] = make([]bool, p.width)
  }

  if tX < 0 || tX >= p.width || tY < 0 || tY >= p.height {
    return errors.New("target is out of bounds")
  }
  if oX == tX && oY == tY {
    return errors.New("target is origin")
  }

  for y, n := range p.nodes {
    for x, _ := range n {
      closedList[x][y] = false
    }
  }

  x := oX
  y := oY
  p.nodes[y][x] = NodeAStar{
    fCost: 0,
    gCost: 0,
    hCost: 0,
    parentX: x,
    parentY: y,
  }

  openNodes := make([]NodeAStar, 0)
  openNodes = append(openNodes, p.nodes[y][x])

  for ; len(openNodes) > 0 && len(openNodes) < p.height*p.width; {
    var node NodeAStar = openNodes[0]
    // remove top node from array
    openNodes = openNodes[1:]

    // Add node to closed array.
    closedList[node.y][node.x] = true

    // Iterate through our eight directions
    for i := -1; i != 1; i+=2 {
      for j := -1; j != 1; j+=2 {
        y = node.y - i
        x = node.x - j
        // Skip past if it is out of bounds.
        if y < 0 || y >= p.width || x < 0 || x >= p.height {
          continue
        }
        // If it is our destination then we've found a path.
        if y == tY && x == tX {
          p.nodes[y][x].parentY = node.y
          p.nodes[y][x].parentX = node.x
          p.tracePath(oX, oY)
          p.foundPath = true
          return nil
        }
        if closedList[y][x] == false && p.nodes[y][x].mCost != math.MaxUint32 {
          g := node.gCost + 1 + float64(node.mCost)
          h := p.calculateH(y, x, tX, tY)
          f := g + h
          //
          if p.nodes[y][x].fCost == math.MaxFloat64 || p.nodes[y][x].fCost > f {
            openNodes = append(openNodes, NodeAStar{
              fCost: f,
              x: x,
              y: y,
            })
            p.nodes[y][x].fCost = f
            p.nodes[y][x].gCost = g
            p.nodes[y][x].hCost = h
            p.nodes[y][x].parentY = y
            p.nodes[y][x].parentX = x
          }
        }
      }
    }
    //
  }
  return errors.New("could not find path")
}

func (p *PathAStar) calculateH(x, y int, tX, tY int) float64 {
  return math.Sqrt(float64((y-tY)*(y-tY) + (x-tX)*(x-tX)))
}

func (p *PathAStar) tracePath(tX, tY int) {
  y := tY
  x := tY
  path := make([]NodeAStar, 0)

  for ; p.nodes[y][x].parentY != y && p.nodes[y][x].parentX != x; {
    path = append(path, NodeAStar{
      x: x,
      y: y,
    })
    tempX := p.nodes[y][x].parentX
    tempY := p.nodes[y][x].parentY
    x = tempX
    y = tempY
  }

  path = append(path, NodeAStar{
    x: x,
    y: y,
  })
  for ; len(path) > 0; {
    node := path[len(path)-1]
    path = path[:len(path)-1]
    fmt.Printf("-> (%d, %d) ", node.y, node.x)
  }
}

func (p *PathAStar) HasRoute() bool {
  return p.foundPath
}

func (p *PathAStar) RouteSize() int {
  return 1
}
