package main

type MaterialType int

const (
	TypeSolid MaterialType = iota
	TypeGas
	TypeFluid
)

type MaterialDef struct {
	ID           string
	Name         string
	Type         MaterialType
	FlowConstant float64
	Density      float64 // kg/m^3 (or arbitrary game units)
	GasConstant  float64 // For gases
}

var (
	Water = MaterialDef{ID: "water", Name: "Water", Type: TypeFluid, FlowConstant: 0.5, Density: 1000.0}
	Steam = MaterialDef{ID: "steam", Name: "Steam", Type: TypeGas, Density: 0.6, GasConstant: 200.0} // density varies, using base
	Coal  = MaterialDef{ID: "coal", Name: "Coal", Type: TypeSolid, Density: 1500.0}
)
