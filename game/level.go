// Copyright 2021 The Ebiten Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
)

// Level represents a Game level.
type Level struct {
	w, h int

	tiles    [][]*Tile // (Y,X) array of tiles
	tileSize int
	entities []*Entity
}

func NewLevel() (*Level, error) {
	l := &Level{
		w:        4,
		h:        4,
		tileSize: 32,
		entities: make([]*Entity, 0),
	}

	ss, err := LoadSpriteSheet(l.tileSize)
	if err != nil {
		return nil, fmt.Errorf("failed to load spritesheet: %s", err)
	}

	l.tiles = make([][]*Tile, l.h)
	for y := 0; y < l.h; y++ {
		l.tiles[y] = make([]*Tile, l.w)
		for x := 0; x < l.w; x++ {
			l.tiles[y][x] = &Tile{}

			floorComp := &Reservoir{
				Basics: Basics{
					Identifier: ".",
					Color:      [3]byte{100, 100, 100},
				},
			}
			floorEntity := NewEntity(x, y, floorComp, ss.Floor, 0)
			l.tiles[y][x].AddEntity(floorEntity)
			l.entities = append(l.entities, floorEntity)
		}
	}

	// Add reservoir at (1,1)
	res1 := &Reservoir{
		Basics: Basics{
			Identifier: "A",
			Color:      [3]byte{0, 0, 255},
		},
		Structurals: Structurals{
			MaxCapacity:     1000,
			CurrentCapacity: 1000,
			Area:            5.0,
		},
	}
	res1Entity := NewEntity(1, 1, res1, ss.Reservoir1, 1)
	l.tiles[1][1].AddEntity(res1Entity)
	l.entities = append(l.entities, res1Entity)

	// Add pipes
	pipe1 := NewPipe(res1, nil, 1.0, 1.0) // we'll fix the connection later
	pipe1Entity := NewEntity(1, 2, pipe1, ss.PipeEnterLeft, 1)
	l.tiles[1][2].AddEntity(pipe1Entity)
	l.entities = append(l.entities, pipe1Entity)

	// Add reservoir at (1,5)
	res2 := &Reservoir{
		Basics: Basics{
			Identifier: "B",
			Color:      [3]byte{0, 0, 255},
		},
		Structurals: Structurals{
			MaxCapacity:     1000,
			CurrentCapacity: 0,
			Area:            5.0,
		},
	}
	res2Entity := NewEntity(1, 5, res2, ss.Reservoir1, 1)
	l.tiles[1][3].AddEntity(res2Entity)
	l.entities = append(l.entities, res2Entity)

	// Wire up the pipe connection
	pipe1.From = res1
	pipe1.To = res2

	return l, nil
}

// Tile returns the tile at the provided coordinates, or nil.
func (l *Level) Tile(x, y int) *Tile {
	if x >= 0 && y >= 0 && x < l.w && y < l.h {
		return l.tiles[y][x]
	}
	return nil
}

// Size returns the size of the Level.
func (l *Level) Size() (width, height int) {
	return l.w, l.h
}
