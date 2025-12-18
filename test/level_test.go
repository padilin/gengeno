package test

import (
	"os"
	"testing"

	"github.com/padilin/gengeno/game"
)

func TestNewLevel(t *testing.T) {
	// Need to handle missing assets gracefully or skip test
	if _, err := os.Stat("assets/floor-tile-1.png"); os.IsNotExist(err) {
		t.Skip("Skipping TestNewLevel because assets/floor-tile-1.png is missing")
	}

	g := &game.Game{System: &game.System{}}
	l, err := game.NewLevel(g)
	if err != nil {
		t.Errorf("NewLevel() error = %v", err)
		return
	}
	if l == nil {
		t.Fatal("NewLevel() returned nil")
	}
	// w, h are likely unexported or we need public accessors. size() returns them but implementation details are hidden.
	// l.w and l.h are not accessible.

	w, h := l.Size()
	if w != 4 || h != 4 {
		t.Errorf("NewLevel() size = %dx%d, want 4x4", w, h)
	}

	// Can't check l.entities or l.System if unexported or no accessors.
	// l.System seems exported, let's check.
	// Checked level.go: System field IS exported in Level struct.
	if l.System != g.System {
		t.Error("NewLevel() System check failed")
	}
}

func TestLevel_Tile(t *testing.T) {
	// t.Skip("Skipping TestLevel_Tile due to setupTestLevel complexity and field access")
	l := setupTestLevel(t)
	// Default level is 4x4
	tile := l.Tile(0, 0)
	if tile == nil {
		t.Error("Tile(0,0) returned nil")
	}
	// Check out of bounds
	if l.Tile(-1, 0) != nil {
		t.Error("Tile(-1,0) should be nil")
	}
	if l.Tile(4, 4) != nil {
		t.Error("Tile(4,4) should be nil")
	}
}

func TestLevel_Size(t *testing.T) {
	// Initialize Level literal with exported fields Width, Height
	l := &game.Level{Width: 10, Height: 20}
	w, h := l.Size()
	if w != 10 || h != 20 {
		t.Errorf("Level{Width: 10, Height: 20}.Size() = %d, %d, want 10, 20", w, h)
	}
}
