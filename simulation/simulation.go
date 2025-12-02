package simulation

import (
	"fmt"
	"math"

	"gengeno/components"
)

// Connection represents the link between 2 components
type Connection struct {
	A        components.Component
	B        components.Component
	FlowArea float64
}

// Simulation holdes the state and logic for simulation.
type Simulation struct {
	Connections []Connection
	AllComps    []components.Component
	Tick        int
	Status      string
}

func NewSimulation(connections []Connection, allComps []components.Component) *Simulation {
	return &Simulation{
		Connections:   connections,
		AllComps:      allComps,
		Status:        "Ready",
	}
}

func (s *Simulation) Update() {
	s.Tick++
	if s.Tick%10 != 0 {
		return
	}

	// Step 1: Determine changes
	for _, conn := range s.Connections {
		s.equalizePressure(conn.A, conn.B, conn.FlowArea)
	}

	// Step 2: Apply pressure changes
	for _, comp := range s.AllComps {
		s.updatePressure(comp)
	}

	s.Status = fmt.Sprintf("Tick %d", s.Tick)
}

func (s *Simulation) updatePressure(c components.Component) {
	s_structs := c.GetStructurals()
	if s_structs != nil && s_structs.Area > 0 {
		s_structs.Pressure = s_structs.CurrentCapacity / s_structs.Area
	}
}

func (s *Simulation) equalizePressure(compA, compB components.Component, flowArea float64) float64 {
	aStructs := compA.GetStructurals()
	bStructs := compB.GetStructurals()

	if aStructs == nil || bStructs == nil {
		return 0
	}

	var highPressureComp, lowPressureComp *components.Structurals
	pressureDifference := aStructs.Pressure - bStructs.Pressure

	if pressureDifference > 0 {
		highPressureComp = aStructs
		lowPressureComp = bStructs
	} else if pressureDifference < 0 {
		highPressureComp = bStructs
		lowPressureComp = aStructs
		pressureDifference = -pressureDifference
	} else {
		return 0
	}

	const flowConstant = 0.5
	flowAmount := math.Sqrt(pressureDifference) * flowArea * flowConstant

	if flowAmount > highPressureComp.CurrentCapacity {
		flowAmount = highPressureComp.CurrentCapacity
	}

	if flowAmount > (lowPressureComp.MaxCapacity - lowPressureComp.CurrentCapacity) {
		flowAmount = lowPressureComp.MaxCapacity - lowPressureComp.CurrentCapacity
	}

	if flowAmount < 0.001 { // Prevent tiny, meaningless transfers
		return 0
	}

	highPressureComp.CurrentCapacity -= flowAmount
	lowPressureComp.CurrentCapacity += flowAmount

	return flowAmount
}

// transferFluid is generic function to calculate movement of fluid between source and destination.
func (s *Simulation) transferFluid(source, destination components.Component, pipeArea float64) float64 {
	sourceStructs := source.GetStructurals()
	destStructs := destination.GetStructurals()

	if sourceStructs == nil || destStructs == nil {
		return 0
	}

	pressureDifference := sourceStructs.Pressure - destStructs.Pressure
	if pressureDifference <= 0 {
		return 0
	}

	const flowConstant = 0.75
	flowAmount := math.Sqrt(pressureDifference) * pipeArea * flowConstant
	if flowAmount > sourceStructs.CurrentCapacity {
		flowAmount = sourceStructs.CurrentCapacity
	}
	if flowAmount > (destStructs.MaxCapacity - destStructs.CurrentCapacity) {
		flowAmount = destStructs.MaxCapacity - destStructs.CurrentCapacity
	}
	if flowAmount <= 0 {
		return 0
	}
	sourceStructs.CurrentCapacity -= flowAmount
	destStructs.CurrentCapacity += flowAmount

	return flowAmount
}

func (s *Simulation) processPipeFlow(pipe *components.Pipe) {
	if len(pipe.Inputs) > 0 {
		inComp := pipe.Inputs[0]
		s.transferFluid(inComp, pipe, pipe.FlowArea)
	}

	if len(pipe.Outputs) > 0 {
		outComp := pipe.Outputs[0]
		s.transferFluid(pipe, outComp, pipe.FlowArea)
	}
}
