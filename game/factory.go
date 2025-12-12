package main

type EntityConfig struct {
	Type       string // "Reservoir", "Pipe", "Wall", "Generator"
	X, Y       int    // Grid coordinates
	Identifier string // Display ID (e.g. "A", "B")

	// Physical Properties
	MaxVolume  float64
	InitialQty float64
	Area       float64
	Contents   *MaterialDef // "Water", "Steam", etc.

	// Component Specifics
	PipeLength float64
	PipeRadius float64

	// Visuals
	Sprite string
}

// Spawn creates an entity based on config and registers it to the level and system.
func (l *Level) Spawn(c EntityConfig) (*Entity, error) {
	var comp Component
	var ent *Entity

	// Defaults
	if c.Area == 0 {
		c.Area = 5.0
	}
	if c.MaxVolume == 0 {
		switch c.Type {
		case "Reservoir":
			c.MaxVolume = 1000.0
		case "Pipe":
			// Calculated from radius/length if not set, but verified later
		}
	}

	// Prepare Contents
	var initialContents []MaterialDef
	if c.InitialQty > 0 {
		if c.Contents != nil {
			initialContents = []MaterialDef{*c.Contents}
		} else {
			initialContents = []MaterialDef{Water} // Default to water
		}
	}

	switch c.Type {
	case "Reservoir":
		res := &Reservoir{
			Basics: Basics{
				Identifier: c.Identifier,
				Color:      [3]byte{0, 0, 255}, // Default blue-ish
			},
			Structurals: Structurals{
				MaxVolume: c.MaxVolume,
				Area:      c.Area,
				Quantity:  c.InitialQty,
				Contents:  initialContents,
			},
		}
		comp = res
		// Create Reservoir Entity (visuals)
		ent = NewReservoirEntity(c.X, c.Y, res, 1)

	case "Pipe":
		// Pipe requires special handling if we want to connect it here,
		// but Spawn might just create the unconnected pipe for now.
		// Connection usually requires references to other components.
		// For now, we creating a pipe component.
		if c.PipeLength == 0 {
			c.PipeLength = 1.0
		}
		if c.PipeRadius == 0 {
			c.PipeRadius = 0.5
		}

		p := NewPipe(nil, nil, c.PipeLength, c.PipeRadius)
		p.Identifier = c.Identifier // Pipe usually doesn't show ID, but for debug
		p.Quantity = c.InitialQty

		comp = p

		sprite := c.Sprite
		if sprite == "" {
			sprite = "pipe_horizontal"
		} // fallback
		ent = NewPipeEntity(c.X, c.Y, p, sprite, 1)
	}

	if ent != nil {
		l.AddEntity(ent)

		// Auto-register to System
		if l.System != nil && comp != nil {
			if pipe, ok := comp.(*Pipe); ok {
				l.System.Pipes = append(l.System.Pipes, pipe)
			} else {
				l.System.Nodes = append(l.System.Nodes, comp)
			}
		}
	}

	return ent, nil
}

// AddEntity handles adding to tiles and internal list
func (l *Level) AddEntity(e *Entity) {
	if l.Tile(e.X, e.Y) != nil {
		l.Tile(e.X, e.Y).AddEntity(e)
		l.entities = append(l.entities, e)
	}
}
