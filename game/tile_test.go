package main

import (
	"testing"
)

func TestTile_AddEntity(t *testing.T) {
	tr := &Tile{}
	ent := &Entity{X: 1, Y: 1} // Pointers

	tr.AddEntity(ent)

	if len(tr.entities) != 1 {
		t.Errorf("Tile.entities len = %d, want 1", len(tr.entities))
	}
	if tr.entities[0] != ent {
		t.Error("Tile.entities[0] does not match added entity")
	}
}

// TestTile_Draw skipped as it involves complex Ebiten image mocking
