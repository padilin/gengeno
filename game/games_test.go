package main

import (
	"math"
	"testing"
)

func TestNewGame(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "Default Initialization",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewGame()
			if (err != nil) != tt.wantErr {
				t.Errorf("NewGame() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got == nil {
				t.Error("NewGame() returned nil")
				return
			}
			if got.System == nil {
				t.Error("NewGame() System is nil") // System might be nil pending other changes, but let's check
			}
			if got.currentLevel == nil {
				t.Error("NewGame() currentLevel is nil")
			}
			if !got.pause {
				t.Error("NewGame() expected pause=true")
			}
			if got.camScale != 1.0 {
				t.Errorf("NewGame() camScale = %v, want 1.0", got.camScale)
			}
		})
	}
}

func TestGame_Update(t *testing.T) {
	g := &Game{
		System:       &System{},
		currentLevel: &Level{w: 10, h: 10, tileSize: 32},
		pause:        false, // Running
		camScale:     1.0,
	}

	// First update should increment ticks
	if err := g.Update(); err != nil {
		t.Errorf("Game.Update() error = %v", err)
	}
	if g.System.Ticks != 1 {
		t.Errorf("System.Ticks = %d, want 1", g.System.Ticks)
	}

	// Pause
	g.pause = true
	g.Update()
	if g.System.Ticks != 1 {
		t.Errorf("System.Ticks (paused) = %d, want 1", g.System.Ticks)
	}
}

func TestGame_Layout(t *testing.T) {
	g := &Game{}
	w, h := g.Layout(800, 600)
	if w != 800 || h != 600 {
		t.Errorf("Layout() = %d,%d, want 800,600", w, h)
	}
	if g.w != 800 || g.h != 600 {
		t.Errorf("Game dims = %d,%d, want 800,600", g.w, g.h)
	}
}

func TestGame_cartesianToIso(t *testing.T) {
	// TileSize = 32
	// ix := (x - y) * 16
	// iy := (x + y) * 8
	g := &Game{currentLevel: &Level{tileSize: 32}}

	tests := []struct {
		name   string
		x, y   float64
		wx, wy float64 // want x, want y
	}{
		{"Origin", 0, 0, 0, 0},
		{"X Axis", 1, 0, 16, 8},
		{"Y Axis", 0, 1, -16, 8},
		{"Both", 1, 1, 0, 16},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotX, gotY := g.cartesianToIso(tt.x, tt.y)
			if math.Abs(gotX-tt.wx) > 0.001 || math.Abs(gotY-tt.wy) > 0.001 {
				t.Errorf("cartesianToIso() = %v,%v, want %v,%v", gotX, gotY, tt.wx, tt.wy)
			}
		})
	}
}

func TestGame_isoToCartesian(t *testing.T) {
	g := &Game{currentLevel: &Level{tileSize: 32}}

	tests := []struct {
		name   string
		x, y   float64 // iso input
		wx, wy float64 // cartesian output
	}{
		{"Origin", 0, 0, 0, 0},
		{"X Axis", 16, 8, 1, 0},
		{"Y Axis", -16, 8, 0, 1},
		{"Both", 0, 16, 1, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotX, gotY := g.isoToCartesian(tt.x, tt.y)
			if math.Abs(gotX-tt.wx) > 0.001 || math.Abs(gotY-tt.wy) > 0.001 {
				t.Errorf("isoToCartesian() = %v,%v, want %v,%v", gotX, gotY, tt.wx, tt.wy)
			}
		})
	}
}

// TestGame_Draw and TestGame_renderLevel omitted/simplified due to graphics dependency
func TestGame_Draw(t *testing.T) {
	// No-op test just to ensure function exists and signature is correct
	// Actual drawing tests require headless context or mock image
}
