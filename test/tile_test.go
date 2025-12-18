package test

import (
	"testing"

	"github.com/padilin/gengeno/game"
)

func TestTile_AddEntity(t *testing.T) {
	tr := &game.Tile{}
	ent := &game.Entity{X: 1, Y: 1} // Pointers

	tr.AddEntity(ent)

	// t.Log("Skipping internal state check for Tile.entities (unexported)")
	if len(tr.Entities()) != 1 {
		t.Errorf("Tile.entities len = %d, want 1", len(tr.Entities()))
	}
	if tr.Entities()[0] != ent {
		t.Error("Tile.entities[0] does not match added entity")
	}
}

// TestTile_Draw skipped as it involves complex Ebiten image mocking
