package test

import (
	"testing"

	"github.com/padilin/gengeno/game"
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
			got, err := game.NewGame()
			if (err != nil) != tt.wantErr {
				t.Errorf("NewGame() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got == nil {
				t.Error("NewGame() returned nil")
				return
			}
			if got.System == nil {
				t.Error("NewGame() System is nil")
			}
			// currentLevel, pause, camScale are unexported and cannot be checked from outside
		})
	}
}

func TestGame_Update(t *testing.T) {
	// Game struct has unexported fields, we rely on NewGame or must interact via public methods.
	g, err := game.NewGame()
	if err != nil {
		t.Fatalf("Failed to create game: %v", err)
	}

	g.SetPause(false) // Unpause to allow ticks

	// First update should increment ticks
	if err := g.Update(); err != nil {
		t.Errorf("Game.Update() error = %v", err)
	}
	if g.System.Ticks != 1 {
		t.Errorf("System.Ticks = %d, want 1", g.System.Ticks)
	}

	// Cannot verify pause logic directly as 'pause' field is unexported
	// But we verified SetPause allows update.
}

func TestGame_Layout(t *testing.T) {
	g, _ := game.NewGame()
	w, h := g.Layout(800, 600)
	if w != 800 || h != 600 {
		t.Errorf("Layout() = %d,%d, want 800,600", w, h)
	}
	// Cannot check g.w, g.h (unexported)
}

func TestGame_CartesianToIso(t *testing.T) {
	// t.Skip("Skipping TestGame_cartesianToIso (unexported method)")
	g, _ := game.NewGame()
	// Mock level settings if needed?
	// The default level has tileSize 32
	// cartToIso(0,0) -> 0,0
	x, y := g.CartesianToIso(0, 0)
	if x != 0 || y != 0 {
		t.Errorf("CartesianToIso(0,0) = %f,%f, want 0,0", x, y)
	}
	// CartesianToIso(1,0) -> (1-0)*16, (1+0)*8 -> 16, 8
	// tileSize/2 = 16, tileSize/4 = 8
	x, y = g.CartesianToIso(1, 0)
	if x != 16 || y != 8 {
		t.Errorf("CartesianToIso(1,0) = %f,%f, want 16,8", x, y)
	}
}

func TestGame_IsoToCartesian(t *testing.T) {
	// t.Skip("Skipping TestGame_isoToCartesian (unexported method)")
	g, _ := game.NewGame()
	// IsoToCartesian(0,0) -> 0,0
	cx, cy := g.IsoToCartesian(0, 0)
	if cx != 0 || cy != 0 {
		t.Errorf("IsoToCartesian(0,0) = %f,%f, want 0,0", cx, cy)
	}
	// IsoToCartesian(16, 8) -> 1, 0
	cx, cy = g.IsoToCartesian(16, 8)
	if cx != 1 || cy != 0 {
		t.Errorf("IsoToCartesian(16,8) = %f,%f, want 1,0", cx, cy)
	}
}

func TestGame_Draw(t *testing.T) {
	// No-op
}
