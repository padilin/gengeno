package components

import (
	"gengeno/materials"
)

// Component is an interface that all game components should implement.
// It provides methods for getting visual information about the component.
type Component interface {
	GetIdentifier() string
	GetColor() (r, g, b byte)
	GetStructurals() *Structurals
}

// Basics contains basic information for a component..
type Basics struct {
	Id         int
	Identifier string
	Color      [3]byte // R, G, B values for the component's color
}

// GetIdentifier returns the single-letter representation of the component.
func (b *Basics) GetIdentifier() string {
	return b.Identifier
}

// GetColor returns the color of the component.
func (b *Basics) GetColor() (byte, byte, byte) {
	return b.Color[0], b.Color[1], b.Color[2]
}

// Structurals contains structural properties for a component.
type Structurals struct {
	MaxCapacity     float64
	CurrentCapacity float64
	MaxHeat         int
	CurrentHeat     int
	MaxPressure     int
	Pressure        float64
	Area            float64
}

// InOut contains properties to connect to other components.
type InOut struct {
	Inputs  []Component
	Outputs []Component
}

// Reservoir represents any component to hold MaterialDef.
type Reservoir struct {
	Basics
	Structurals
	InOut
	Contents materials.MaterialDef
}

func (r *Reservoir) GetStructurals() *Structurals {
	return &r.Structurals
}

// PipeSystem is a container for Pipe components.
type PipeSystem struct {
	Basics
	Pipes []Pipe
}

// Pipe moves MaterialDef between components.
type Pipe struct {
	Basics
	Structurals
	InOut
	Contents materials.MaterialDef
	FlowArea float64
}

func (p *Pipe) GetStructurals() *Structurals {
	return &p.Structurals
}

// Generator represents any component to generate.
type Generator struct {
	Basics
	Structurals
	InOut
	Contents []materials.MaterialDef
}
