package components

import (
	"gengeno/materials"
)

type basics struct {
	id int
	identifier string
}

type structurals struct {
	basics
	maxCapacity int
	currentCapacity int
	maxHeat int
	currentHeat int
	maxPressure int
	currentPressure int
}

type Reservoir struct {
	basics
	structurals
	contents materials.MaterialDef
}

type PipeSystem struct {
	basics
	pipes []Pipe
}

type Pipe struct {
	basics
	structurals
	contents materials.MaterialDef
}

type Generator struct {
	basics
	structurals
	contents materials.MaterialDef
}
