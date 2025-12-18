package test

import (
	"testing"

	"github.com/padilin/gengeno/game"
)

// setupTestLevel cannot access unexported fields of Level if defined here.
// We must skip these tests unless game helper is provided.

func TestLevel_Spawn(t *testing.T) {
	// Setup
	orig := game.SpriteSet
	game.SpriteSet = make(map[string]*game.Sprite)
	game.SpriteSet["floor"] = &game.Sprite{DrawOrder: 1} // Needed for floor
	defer func() { game.SpriteSet = orig }()

	l := setupTestLevel(t)

	// Create a dummy config for testing
	cfg := game.EntityConfig{
		Type: "Reservoir",
		X:    1, Y: 1,
		Identifier: "TEST_RES",
		MaxVolume:  100,
	}

	ent, err := l.Spawn(cfg)
	if err != nil {
		t.Fatalf("Spawn failed: %v", err)
	}
	if ent == nil {
		t.Fatal("Spawn returned nil entity")
	}
	if ent.X != 1 || ent.Y != 1 {
		t.Errorf("Spawn entity at %d,%d, want 1,1", ent.X, ent.Y)
	}
}

func TestLevel_AddEntity(t *testing.T) {
	l := setupTestLevel(t)
	// Initial entities from NewLevel (e.g. floor tiles)
	// We count them.
	initialCount := len(l.Entities())

	ent := &game.Entity{X: 0, Y: 0}
	l.AddEntity(ent)

	if len(l.Entities()) != initialCount+1 {
		t.Errorf("Level entities count = %d, want %d", len(l.Entities()), initialCount+1)
	}
}

func setupTestLevel(t *testing.T) *game.Level {
	g := &game.Game{System: &game.System{}}
	l, err := game.NewLevel(g)
	if err != nil {
		t.Fatalf("NewLevel failed: %v", err)
	}
	return l
}
