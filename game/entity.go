package main

type SpriteSelector func(e *Entity) *Sprite

type Entity struct {
	X, Y      int
	Component Component
	Selector  SpriteSelector
	Sprite    *Sprite
}

func (e *Entity) CurrentSprite() *Sprite {
	if e == nil {
		return nil
	}
	if e.Selector != nil {
		return e.Selector(e)
	}
	return e.Sprite // fallback to static sprite
}

// === Generic Selector Factory ===
// StaticSpriteSelector returns a selector that always returns the same sprite key.
func StaticSpriteSelector(key string) SpriteSelector {
	return func(e *Entity) *Sprite {
		return spriteSet[key]
	}
}

// FillPercentSelector returns a selector that picks sprites based on fill %.
// Useful for tanks, pipes, etc. that have capacity-based states.
func FillPercentSelector(stateMap map[string]string) SpriteSelector {
	return func(e *Entity) *Sprite {
		s := e.Component.GetStructurals()
		if s == nil {
			return nil
		}

		p := 0.0
		if s.MaxVolume > 0 {
			p = s.Volume / s.MaxVolume
		}

		// Find the appropriate state key based on fill %
		var state string
		switch {
		case p > 0.9:
			state = "full"
		case p > 0.6:
			state = "high"
		case p > 0.3:
			state = "mid"
		case p > 0.0:
			state = "low"
		default:
			state = "empty"
		}

		// Look up the sprite in the provided state map
		if key, ok := stateMap[state]; ok {
			return spriteSet[key]
		}
		return nil
	}
}

// === Entity Factory Functions ===
// NewEntity creates a simple entity with a static sprite and optional selector.
func NewEntity(x, y int, comp Component, selector SpriteSelector, drawOrder int) *Entity {
	e := &Entity{
		X:         x,
		Y:         y,
		Component: comp,
		Selector:  selector,
	}

	// Set a fallback sprite if selector fails
	if selector != nil {
		e.Sprite = &Sprite{
			Image:     nil, // will use selector
			DrawOrder: drawOrder,
		}
	}

	return e
}

// NewReservoirEntity creates a reservoir entity with fill-based sprite selection.
func NewReservoirEntity(x, y int, comp Component, drawOrder int) *Entity {
	stateMap := map[string]string{
		"full":  "reservoir_full",
		"high":  "reservoir_high",
		"mid":   "reservoir_mid",
		"low":   "reservoir_low",
		"empty": "reservoir_empty",
	}

	return NewEntity(x, y, comp, FillPercentSelector(stateMap), drawOrder)
}

// NewFloorEntity creates a floor tile entity.
func NewFloorEntity(x, y int, comp Component, drawOrder int) *Entity {
	return NewEntity(x, y, comp, StaticSpriteSelector("floor"), drawOrder)
}

// NewPipeEntity creates a pipe entity.
func NewPipeEntity(x, y int, comp Component, spriteKey string, drawOrder int) *Entity {
	return NewEntity(x, y, comp, StaticSpriteSelector(spriteKey), drawOrder)
}
