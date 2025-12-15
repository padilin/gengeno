package main

import (
	"testing"
)

// Mock Tile for testing
func (l *Level) ResetEntities() {
	l.entities = []*Entity{}
	// Reset tiles for testing isolation
	for y := 0; y < l.h; y++ {
		for x := 0; x < l.w; x++ {
			l.tiles[y][x] = &Tile{}
		}
	}
}

func setupTestLevel() *Level {
	l := &Level{
		w:        4,
		h:        4,
		tileSize: 32,
		System:   &System{},
	}
	l.tiles = make([][]*Tile, l.h)
	for y := 0; y < l.h; y++ {
		l.tiles[y] = make([]*Tile, l.w)
		for x := 0; x < l.w; x++ {
			l.tiles[y][x] = &Tile{}
		}
	}
	return l
}

func TestLevel_Spawn(t *testing.T) {
	origSpriteSet := spriteSet
	spriteSet = make(map[string]*Sprite)
	spriteSet["reservoir_full"] = &Sprite{DrawOrder: 5}
	spriteSet["pipe_horizontal"] = &Sprite{DrawOrder: 5}
	defer func() { spriteSet = origSpriteSet }()

	l := setupTestLevel()

	t.Run("Spawn Reservoir", func(t *testing.T) {
		conf := EntityConfig{
			Type:       "Reservoir",
			X:          1,
			Y:          1,
			Identifier: "R1",
			InitialQty: 100,
		}
		ent, err := l.Spawn(conf)
		if err != nil {
			t.Fatalf("Spawn failed: %v", err)
		}
		if ent == nil {
			t.Fatal("Spawn returned nil entity")
		}
		if _, ok := ent.Component.(*Reservoir); !ok {
			t.Error("Spawned entity component is not Reservoir")
		}
		// Check System registration
		if len(l.System.Nodes) == 0 {
			t.Error("System.Nodes is empty, expected registration")
		}
	})

	t.Run("Spawn Pipe", func(t *testing.T) {
		conf := EntityConfig{
			Type:       "Pipe",
			X:          2,
			Y:          2,
			Identifier: "P1",
			PipeLength: 5,
		}
		ent, err := l.Spawn(conf)
		if err != nil {
			t.Fatalf("Spawn failed: %v", err)
		}
		if ent == nil {
			t.Fatal("Spawn returned nil entity")
		}
		p, ok := ent.Component.(*Pipe)
		if !ok {
			t.Error("Spawned entity component is not Pipe")
		} else {
			if p.Length != 5 {
				t.Errorf("Pipe length = %v, want 5", p.Length)
			}
		}
		if len(l.System.Pipes) == 0 {
			t.Error("System.Pipes is empty, expected registration")
		}
	})
}

func TestLevel_AddEntity(t *testing.T) {
	l := setupTestLevel()
	ent := &Entity{X: 0, Y: 0}

	l.AddEntity(ent)

	if len(l.entities) != 1 {
		t.Errorf("Level.entities len = %d, want 1", len(l.entities))
	}
	// Check tile
	// tile := l.Tile(0, 0)
	// Tile implementation detail: AddEntity needs to be verified on Tile itself or via side effect?
	// Actually Level.AddEntity calls Tile.AddEntity. We'll verify Tile tests separately.
	// But we can check if it calls it safely.
	if l.Tile(0, 0) == nil {
		t.Error("Tile(0,0) is nil")
	}
}
