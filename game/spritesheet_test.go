package main

import (
	_ "image/png"
	"os"
	"testing"
)

func TestLoadSpriteSheet(t *testing.T) {
	if _, err := os.Stat("assets/floor-tile-1.png"); os.IsNotExist(err) {
		t.Skip("Skipping TestLoadSpriteSheet because assets/floor-tile-1.png is missing")
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
			got, err := LoadSpriteSheet(tt.tileSize)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadSpriteSheet() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got == nil && !tt.wantErr {
				t.Error("LoadSpriteSheet() returned nil")
			}
			// Verify map is populated
			if len(spriteSet) == 0 {
				t.Error("spriteSet is empty")
			}
		})
	}
}
