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
	Ticks       int
	Status      string
}

func NewSimulation(connections []Connection, allComps []components.Component) *Simulation {
	return &Simulation{
		Connections: connections,
		AllComps:    allComps,
		Status:      "Ready",
		Ticks:       0,
	}
}

func (s *Simulation) Update() {
	s.Ticks++
	if s.Ticks%10 != 0 {
		return
	}

	// // Step 1: Determine changes
	// for _, conn := range s.Connections {
	// 	s.equalizePressure(conn.A, conn.B, conn.FlowArea)
	// }

	// // Step 2: Apply pressure changes
	// for _, comp := range s.AllComps {
	// 	s.updatePressure(comp)
	// }

	s.Status = fmt.Sprintf("Tick %d", s.Ticks)
}

// func (s *Simulation) updatePressure(c components.Component) {
// 	s_structs := c.GetStructurals()
// 	if s_structs != nil && s_structs.Area > 0 {
// 		s_structs.Pressure = s_structs.CurrentCapacity / s_structs.Area
// 	}
// }

// func (s *Simulation) equalizePressure(compA, compB components.Component, flowArea float64) float64 {
// 	aStructs := compA.GetStructurals()
// 	bStructs := compB.GetStructurals()

// 	if aStructs == nil || bStructs == nil {
// 		return 0
// 	}

// 	var highPressureComp, lowPressureComp *components.Structurals
// 	pressureDifference := aStructs.Pressure - bStructs.Pressure

// 	if pressureDifference > 0 {
// 		highPressureComp = aStructs
// 		lowPressureComp = bStructs
// 	} else if pressureDifference < 0 {
// 		highPressureComp = bStructs
// 		lowPressureComp = aStructs
// 		pressureDifference = -pressureDifference
// 	} else {
// 		return 0
// 	}

// 	const flowConstant = 0.5
// 	flowAmount := math.Sqrt(pressureDifference) * flowArea * flowConstant

// 	if flowAmount > highPressureComp.CurrentCapacity {
// 		flowAmount = highPressureComp.CurrentCapacity
// 	}

// 	if flowAmount > (lowPressureComp.MaxCapacity - lowPressureComp.CurrentCapacity) {
// 		flowAmount = lowPressureComp.MaxCapacity - lowPressureComp.CurrentCapacity
// 	}

// 	if flowAmount < 0.001 { // Prevent tiny, meaningless transfers
// 		return 0
// 	}

// 	highPressureComp.CurrentCapacity -= flowAmount
// 	lowPressureComp.CurrentCapacity += flowAmount

// 	return flowAmount
// }

// // transferFluid is generic function to calculate movement of fluid between source and destination.
// func (s *Simulation) transferFluid(source, destination components.Component, pipeArea float64) float64 {
// 	sourceStructs := source.GetStructurals()
// 	destStructs := destination.GetStructurals()

// 	if sourceStructs == nil || destStructs == nil {
// 		return 0
// 	}

// 	pressureDifference := sourceStructs.Pressure - destStructs.Pressure
// 	if pressureDifference <= 0 {
// 		return 0
// 	}

// 	const flowConstant = 0.75
// 	flowAmount := math.Sqrt(pressureDifference) * pipeArea * flowConstant
// 	if flowAmount > sourceStructs.CurrentCapacity {
// 		flowAmount = sourceStructs.CurrentCapacity
// 	}
// 	if flowAmount > (destStructs.MaxCapacity - destStructs.CurrentCapacity) {
// 		flowAmount = destStructs.MaxCapacity - destStructs.CurrentCapacity
// 	}
// 	if flowAmount <= 0 {
// 		return 0
// 	}
// 	sourceStructs.CurrentCapacity -= flowAmount
// 	destStructs.CurrentCapacity += flowAmount

// 	return flowAmount
// }

// func (s *Simulation) processPipeFlow(pipe *components.Pipe) {
// 	if len(pipe.Inputs) > 0 {
// 		inComp := pipe.Inputs[0]
// 		s.transferFluid(inComp, pipe, pipe.FlowArea)
// 	}

// 	if len(pipe.Outputs) > 0 {
// 		outComp := pipe.Outputs[0]
// 		s.transferFluid(pipe, outComp, pipe.FlowArea)
// 	}
// }

// New Simulation
const (
	Gravity      = 9.81
	TimeStep     = 1.0 / 6.0
	FrictionFact = 0.02
	MinorLoss    = 1.5
)

func TotalHead(c components.Component) float64 {
	if c == nil {
		return 0
	}
	s := c.GetStructurals()
	if s == nil {
		return 0
	}
	if s.Area == 0 {
		return s.BaseElevation
	}
	return (s.Volume / s.Area) + s.BaseElevation
}

func ApplyPending(c components.Component) {
	if c == nil {
		return
	}
	r := c.GetStructurals()
	if r == nil {
		return
	}
	r.Volume += r.PendingChange
	r.PendingChange = 0
	if r.Volume < 0 {
		r.Volume = 0
	}
}

type System struct {
	Nodes  []components.Component
	Pipes  []*components.Pipe
	Ticks  int
	Status string
}

func (s *System) Tick() {
	s.Ticks++
	if s.Ticks%10 != 0 {
		return
	}

	for _, p := range s.Pipes {
		// Find difference of head to find flow.
		// Optionally adds in pumphead.
		headFrom := TotalHead(p.From)
		headTo := TotalHead(p.To)
		deltaH := (headFrom + p.PumpHead) - headTo
		direction := 1.0
		if deltaH < 0 {
			direction = -1.0
			deltaH = -deltaH
		}

		// Close enough, stop processing.
		if deltaH < 0.00001 {
			continue
		}

		// Bernoulli
		frictionLoss := FrictionFact * (p.Length / p.Diameter)
		velocity := math.Sqrt((2 * Gravity * deltaH) / (1 + frictionLoss + MinorLoss))

		// Volume to move with direction
		flowVol := velocity * p.Area * TimeStep
		amountMoving := flowVol * direction

		var source, dest components.Component
		if amountMoving > 0 {
			source, dest = p.From, p.To
		} else {
			source, dest = p.To, p.From
			amountMoving = -amountMoving
		}

		srcS := source.GetStructurals()
		dstS := dest.GetStructurals()

		// Pipe fill logic
		spaceInPipe := p.MaxVolume - p.Volume
		if spaceInPipe > 0 {
			if amountMoving <= spaceInPipe {
				srcS.PendingChange -= amountMoving
				p.Volume += amountMoving
			} else {
				srcS.PendingChange -= amountMoving
				p.Volume = p.MaxVolume
				dstS.PendingChange += (amountMoving - spaceInPipe)
			}
		} else {
			srcS.PendingChange -= amountMoving
			dstS.PendingChange += amountMoving
		}
	}

	// Update all Nodes
	for _, node := range s.Nodes {
		ApplyPending(node)
	}

	status := fmt.Sprintf("Tick %d\n", s.Ticks)
	for _, n := range s.Nodes {
		if st := n.GetStructurals(); st != nil {
			status += fmt.Sprintf("  [%s] Vol: %.2f  Head: %.2f\n", n.GetIdentifier(), st.Volume, TotalHead(n))
		} else {
			status += fmt.Sprintf("  [%s] (no structurals)\n", n.GetIdentifier())
		}
	}
	for _, p := range s.Pipes {
		fromID := "?"
		toID := "?"
		if p.From != nil {
			fromID = p.From.GetIdentifier()
		}
		if p.To != nil {
			toID = p.To.GetIdentifier()
		}
		status += fmt.Sprintf("  (Pipe %s -> %s) Vol: %.2f\n", fromID, toID, p.Volume)
	}

	status += "--------------------------------\n"

	s.Status = status
}
