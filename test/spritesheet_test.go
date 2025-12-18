package test

import (
	_ "image/png"
	"os"
	"testing"

	"github.com/padilin/gengeno/game"
)

func TestLoadSpriteSheet(t *testing.T) {
	// Asset should be present in test/assets/floor-tile-1.png or relative path
	_, err := os.Stat("assets/floor-tile-1.png")
	if os.IsNotExist(err) {
		t.Fatalf("Asset assets/floor-tile-1.png missing: %v", err)
	}

	tests := []struct {
		name     string
		tileSize int
		wantErr  bool
	}{
		{
			name:     "Valid Load",
			tileSize: 32,
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := game.LoadSpriteSheet(tt.tileSize)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadSpriteSheet() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got == nil && !tt.wantErr {
				t.Error("LoadSpriteSheet() returned nil")
			}
			// Verify map is populated
			// spriteSet is unexported. Cannot check.
		})
	}
}
