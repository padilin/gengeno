package main

import (
	"os"
	"testing"
)

func TestNewLevel(t *testing.T) {
	// Need to handle missing assets gracefully or skip test
	if _, err := os.Stat("assets/floor-tile-1.png"); os.IsNotExist(err) {
		t.Skip("Skipping TestNewLevel because assets/floor-tile-1.png is missing")
	}

	g := &Game{System: &System{}}
	l, err := NewLevel(g)
	if err != nil {
		t.Errorf("NewLevel() error = %v", err)
		return
	}
	if l == nil {
		t.Fatal("NewLevel() returned nil")
	}
	if l.w != 4 || l.h != 4 {
		t.Errorf("NewLevel() size = %dx%d, want 4x4", l.w, l.h)
	}
	if len(l.entities) == 0 {
		t.Error("NewLevel() entities empty")
	}
	if l.System != g.System {
		t.Error("NewLevel() System check failed")
	}
}

func TestLevel_Tile(t *testing.T) {
	l := setupTestLevel() // Reuse from factory_test.go if package matches, otherwise redefine
	// Since both are package main, they share scope.

	tests := []struct {
		name string
		x    int
		y    int
		want *Tile
	}{
		{"Valid Tile", 0, 0, l.tiles[0][0]},
		{"Out of Bounds X", -1, 0, nil},
		{"Out of Bounds Y", 0, 5, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := l.Tile(tt.x, tt.y)
			if got != tt.want {
				t.Errorf("Level.Tile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLevel_Size(t *testing.T) {
	l := &Level{w: 10, h: 20}
	w, h := l.Size()
	if w != 10 || h != 20 {
		t.Errorf("Level.Size() = %d,%d, want 10,20", w, h)
	}
}
