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
	FlowConstant float64
	Volume int
}

var (
	Water = &MaterialDef{ID: "water", Name: "Water", Type: TypeFluid, FlowConstant: 0.5}
	Steam = &MaterialDef{ID: "steam", Name: "Steam", Type: TypeFluid}
	Coal  = &MaterialDef{ID: "coal", Name: "Coal", Type: TypeSolid}
)
