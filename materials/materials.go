package materials

import ()

type MaterialType int

const (
	TypeSolid MaterialType = iota
	TypeFluid
)

type MaterialDef struct {
	ID string
	Name string
	Type MaterialType
}

var (
	Water = &MaterialDef{ID: "water", Name: "Water", Type: TypeFluid}
	Steam = &MaterialDef{ID: "steam", Name: "Steam", Type: TypeFluid}
	Coal = &MaterialDef{ID: "coal", Name: "Coal", Type: TypeSolid}
)
