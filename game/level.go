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
	System   *System
}

func NewLevel(g *Game) (*Level, error) {
	l := &Level{
		w:        4,
		h:        4,
		tileSize: 32,
		entities: make([]*Entity, 0),
	}

	_, err := LoadSpriteSheet(l.tileSize)
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
			floorEntity := NewFloorEntity(x, y, floorComp, 0)
			l.tiles[y][x].AddEntity(floorEntity)
			l.entities = append(l.entities, floorEntity)
		}
	}

	if g.System == nil {
		g.System = &System{}
	}
	l.System = g.System

	// Add reservoir A
	entA, _ := l.Spawn(EntityConfig{
		Type: "Reservoir",
		X:    1, Y: 1,
		Identifier: "A",
		MaxVolume:  2000,
		InitialQty: 2000,
		Contents:   &Water,
	})

	// Add Pipe
	entPipe, _ := l.Spawn(EntityConfig{
		Type: "Pipe",
		X:    1, Y: 2,
		Identifier: "P1",
		InitialQty: 0.5, // partial filled
		PipeLength: 15.0,
		PipeRadius: 0.5,
		Sprite:     "pipe_enter_left",
	})

	// Add reservoir B
	entB, _ := l.Spawn(EntityConfig{
		Type: "Reservoir",
		X:    1, Y: 3,
		Identifier: "B",
		MaxVolume:  1500,
		InitialQty: 0,
		Contents:   &Water,
	})

	// Wire up the pipe connection
	// We need to extract the components from the entities
	res1 := entA.Component
	res2 := entB.Component
	pipe1 := entPipe.Component.(*Pipe)

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
