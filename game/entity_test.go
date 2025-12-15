package main

import (
	"testing"
)

// Mock implementation of Component for testing
type mockComponent struct {
	structurals Structurals
}

func (m *mockComponent) GetIdentifier() string        { return "M" }
func (m *mockComponent) GetColor() (byte, byte, byte) { return 0, 0, 0 }
func (m *mockComponent) GetStructurals() *Structurals { return &m.structurals }

func TestEntity_CurrentSprite(t *testing.T) {
	// Setup mock sprite set for testing
	originalSpriteSet := spriteSet
	spriteSet = make(map[string]*Sprite)
	spriteSet["test_sprite"] = &Sprite{DrawOrder: 1}
	defer func() { spriteSet = originalSpriteSet }() // Restore after test

	e := &Entity{
		Sprite: &Sprite{DrawOrder: 2},
	}

	t.Run("Fallback to Static Sprite", func(t *testing.T) {
		got := e.CurrentSprite()
		if got == nil || got.DrawOrder != 2 {
			t.Errorf("CurrentSprite() = %v, want matching static sprite", got)
		}
	})

	e.Selector = func(e *Entity) *Sprite {
		return spriteSet["test_sprite"]
	}

	t.Run("Use Selector", func(t *testing.T) {
		got := e.CurrentSprite()
		if got == nil || got.DrawOrder != 1 {
			t.Errorf("CurrentSprite() = %v, want matching test_sprite", got)
		}
	})
}

func TestStaticSpriteSelector(t *testing.T) {
	originalSpriteSet := spriteSet
	spriteSet = make(map[string]*Sprite)
	testSprite := &Sprite{DrawOrder: 5}
	spriteSet["static_key"] = testSprite
	defer func() { spriteSet = originalSpriteSet }()

	selector := StaticSpriteSelector("static_key")
	got := selector(nil)
	if got != testSprite {
		t.Errorf("StaticSpriteSelector() returned wrong sprite")
	}
}

func TestFillPercentSelector(t *testing.T) {
	originalSpriteSet := spriteSet
	spriteSet = make(map[string]*Sprite)
	spriteSet["full"] = &Sprite{DrawOrder: 100}
	spriteSet["empty"] = &Sprite{DrawOrder: 0}
	defer func() { spriteSet = originalSpriteSet }()

	stateMap := map[string]string{
		"full":  "full",
		"empty": "empty",
	}
	selector := FillPercentSelector(stateMap)

	tests := []struct {
		name string
		qty  float64
		max  float64
		want *Sprite
	}{
		{"Full", 100, 100, spriteSet["full"]},
		{"Empty", 0, 100, spriteSet["empty"]},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			comp := &mockComponent{structurals: Structurals{Quantity: tt.qty, MaxVolume: tt.max}}
			e := &Entity{Component: comp}
			got := selector(e)
			if got != tt.want {
				t.Errorf("FillPercentSelector() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewEntity(t *testing.T) {
	comp := &mockComponent{}
	e := NewEntity(10, 20, comp, nil, 5)

	if e.X != 10 || e.Y != 20 {
		t.Errorf("NewEntity() coords = %d,%d, want 10,20", e.X, e.Y)
	}
	if e.Component != comp {
		t.Error("NewEntity() component mismatch")
	}
}

func TestNewReservoirEntity(t *testing.T) {
	comp := &mockComponent{}
	e := NewReservoirEntity(0, 0, comp, 1)
	if e.Selector == nil {
		t.Error("NewReservoirEntity() selector is nil")
	}
}

func TestNewFloorEntity(t *testing.T) {
	comp := &mockComponent{}
	e := NewFloorEntity(0, 0, comp, 1)
	if e.Selector == nil {
		t.Error("NewFloorEntity() selector is nil")
	}
}

func TestNewPipeEntity(t *testing.T) {
	comp := &mockComponent{}
	e := NewPipeEntity(0, 0, comp, "pipe_h", 1)
	if e.Selector == nil {
		t.Error("NewPipeEntity() selector is nil")
	}
}
