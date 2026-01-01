package test

import (
	"testing"

	"github.com/padilin/gengeno/game"
)

// Mock implementation of Component for testing
type mockComponent struct {
	structurals game.Structurals
}

func (m *mockComponent) GetIdentifier() string             { return "M" }
func (m *mockComponent) GetColor() (byte, byte, byte)      { return 0, 0, 0 }
func (m *mockComponent) GetStructurals() *game.Structurals { return &m.structurals }

func TestEntity_CurrentSprite(t *testing.T) {
	// Setup mock sprite set for testing
	// We use the exported game.SpriteSet
	game.SpriteSet = make(map[string]*game.Sprite)
	game.SpriteSet["test"] = &game.Sprite{DrawOrder: 1}

	e := &game.Entity{
		Selector: func(e *game.Entity) *game.Sprite {
			return game.SpriteSet["test"]
		},
	}

	if s := e.CurrentSprite(); s == nil {
		t.Error("CurrentSprite() returned nil")
	} else if s.DrawOrder != 1 {
		t.Errorf("CurrentSprite() order = %d, want 1", s.DrawOrder)
	}
}

func TestStaticSpriteSelector(t *testing.T) {
	game.SpriteSet = make(map[string]*game.Sprite)
	game.SpriteSet["fixed"] = &game.Sprite{DrawOrder: 99}

	sel := game.StaticSpriteSelector("fixed")
	e := &game.Entity{}
	s := sel(e)
	if s == nil {
		t.Fatal("StaticSpriteSelector returned nil")
	}
	if s.DrawOrder != 99 {
		t.Errorf("StaticSpriteSelector order = %d, want 99", s.DrawOrder)
	}
}

func TestFillPercentSelector(t *testing.T) {
	game.SpriteSet = make(map[string]*game.Sprite)
	game.SpriteSet["low"] = &game.Sprite{DrawOrder: 10}
	game.SpriteSet["high"] = &game.Sprite{DrawOrder: 20}

	stateMap := map[string]string{
		"low":  "low",
		"high": "high",
	}
	sel := game.FillPercentSelector(stateMap)

	// Case 1: Low fill
	e := &game.Entity{Component: &mockComponent{structurals: game.Structurals{MaxVolume: 100, Quantity: 10}}}
	s := sel(e)
	if s == nil || s.DrawOrder != 10 {
		t.Error("FillPercentSelector (low) failed")
	}

	// Case 2: High fill
	e.Component = &mockComponent{structurals: game.Structurals{MaxVolume: 100, Quantity: 80}}
	s = sel(e)
	if s == nil || s.DrawOrder != 20 {
		t.Error("FillPercentSelector (high) failed")
	}
}

func TestNewEntity(t *testing.T) {
	comp := &mockComponent{}
	e := game.NewEntity(10, 20, comp, nil, 5)

	if e.X != 10 || e.Y != 20 {
		t.Errorf("NewEntity() coords = %d,%d, want 10,20", e.X, e.Y)
	}
	if e.Component != comp {
		t.Error("NewEntity() component mismatch")
	}
}

func TestNewReservoirEntity(t *testing.T) {
	comp := &mockComponent{}
	e := game.NewReservoirEntity(0, 0, comp, 1)
	if e.Selector == nil {
		t.Error("NewReservoirEntity() selector is nil")
	}
}

func TestNewFloorEntity(t *testing.T) {
	comp := &mockComponent{}
	e := game.NewFloorEntity(0, 0, comp, 1)
	if e.Selector == nil {
		t.Error("NewFloorEntity() selector is nil")
	}
}

func TestNewPipeEntity(t *testing.T) {
	comp := &mockComponent{}
	e := game.NewPipeEntity(0, 0, comp, "pipe_h", 1)
	if e.Selector == nil {
		t.Error("NewPipeEntity() selector is nil")
	}
}
