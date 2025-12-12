package main

import (
	"log"
	"math"
)

const (
	Gravity      = 9.81
	TimeStep     = 1.0 / 6.0
	FrictionFact = 0.02
	MinorLoss    = 1.5
)

func GetMaterial(c Component) *MaterialDef {
	s := c.GetStructurals()
	if s == nil || len(s.Contents) == 0 {
		return &Water // Default to Water if empty
	}
	return &s.Contents[0]
}

func TotalHead(c Component) float64 {
	if c == nil {
		return 0
	}
	s := c.GetStructurals()
	if s == nil {
		return 0
	}

	mat := GetMaterial(c)

	// Case 1: Fluid (Hydrostatic)
	if mat.Type == TypeFluid {
		if s.Area == 0 {
			return s.BaseElevation
		}
		// Volume = Mass / Density
		// Head = Volume / Area = (Mass / Density) / Area
		if mat.Density == 0 {
			return s.BaseElevation
		} // prevent div/0
		return (s.Quantity / (mat.Density * s.Area)) + s.BaseElevation
	}

	// Case 2: Gas (Compressible)
	// Head = P / (rho_water * g)
	// P = (Mass * R) / Vol_container
	if mat.Type == TypeGas {
		if mat.GasConstant > 0 {
			if s.MaxVolume == 0 {
				return 0
			}
			pressure := (s.Quantity * mat.GasConstant) / s.MaxVolume
			// Convert pressure to head (meters of water)
			// P = rho * g * h => h = P / (rho * g)
			// Using Water density as reference for Head
			return pressure / (Water.Density * Gravity)
		}
	}

	return s.BaseElevation
}

func ApplyPending(c Component) {
	if c == nil {
		return
	}
	r := c.GetStructurals()
	if r == nil {
		return
	}
	log.Printf("ApplyPending %v quantity=%.2f change=%.3f", identifier(c), r.Quantity, r.PendingChange)
	r.Quantity += r.PendingChange
	r.PendingChange = 0
	if r.Quantity < 0 {
		r.Quantity = 0
	}
}

type System struct {
	Nodes []Component
	Pipes []*Pipe
	Ticks int
}

func (s *System) Tick() {
	s.Ticks++
	if s.Ticks%10 != 0 {
		return
	}

	log.Printf("SIM Ticks=%d Nodes=%d Pipes=%d", s.Ticks, len(s.Nodes), len(s.Pipes))
	for i, p := range s.Pipes {
		in := p.From
		out := p.To
		var inS, outS *Structurals
		if in != nil {
			inS = in.GetStructurals()
		}
		if out != nil {
			outS = out.GetStructurals()
		}
		log.Printf("Pipe[%d] area=%.3f len=%.3f from=%v quantity=%.2f pres=%.3f -> to=%v quantity=%.2f pres=%.3f",
			i, p.Area, p.Length,
			identifier(in), inS.Quantity, pres(inS),
			identifier(out), outS.Quantity, pres(outS))
	}

	for _, p := range s.Pipes {
		// Calculate Flow 1: Source -> Pipe
		if p.From != nil {
			calculateFlow(p.From, p, p.PumpHead) // PumpHead usually applies to flow *through* pipe?
			// Let's assume PumpHead helps move From -> To.
			// For simplicity: From -> Pipe (Gravity/Pressure), Pipe -> To (Gravity/Pressure + Pump?)
			// Or apply PumpHead to the whole path?
			// Let's stick to simple Head diff for Input.
		}

		// Calculate Flow 2: Pipe -> Destination
		if p.To != nil {
			calculateFlow(p, p.To, p.PumpHead)
		}
	}

	// Update all Nodes and Pipes
	for _, node := range s.Nodes {
		ApplyPending(node)
	}
	for _, pipe := range s.Pipes {
		ApplyPending(pipe)
	}
}

func calculateFlow(from, to Component, pumpHead float64) {
	if from == nil || to == nil {
		return
	}

	pFrom := from.GetStructurals()
	pTo := to.GetStructurals()

	headFrom := TotalHead(from)
	headTo := TotalHead(to)

	deltaH := (headFrom + pumpHead) - headTo

	if deltaH < 0.00001 {
		return
	}

	// Simple flow calc
	// Assuming generic connection properties?
	// If one is a Pipe, use its geometry?
	// We need 'Connection Geometry'.
	// For Source->Pipe, use Pipe geometry.
	// For Pipe->Dest, use Pipe geometry.

	var area, length, radius float64
	if pipe, ok := to.(*Pipe); ok {
		area = pipe.Area
		length = pipe.Length / 2 // Half length for input?
		radius = pipe.Radius
	} else if pipe, ok := from.(*Pipe); ok {
		area = pipe.Area
		length = pipe.Length / 2 // Half length for output?
		radius = pipe.Radius
	} else {
		// Fallback
		area = 1.0
		length = 1.0
		radius = 0.5
	}

	// Bernoulli
	frictionLoss := FrictionFact * (length / (2 * radius))
	velocity := math.Sqrt((2 * Gravity * deltaH) / (1 + frictionLoss + MinorLoss))

	flowVol := velocity * area * TimeStep

	// Check content/density
	mat := GetMaterial(from)
	density := mat.Density
	if density == 0 {
		density = 1000
	}

	amountMoving := flowVol * density

	// Constraint: Source Quantity
	if amountMoving > pFrom.Quantity {
		amountMoving = pFrom.Quantity
	}

	// Constraint: Dest Capacity
	// If Dest is Pipe or Reservoir, check MaxVolume
	currentVol := pTo.Quantity / density
	spaceVol := pTo.MaxVolume - currentVol

	if spaceVol <= 0 {
		return // Full
	}

	volMoving := amountMoving / density
	if volMoving > spaceVol {
		amountMoving = spaceVol * density
	}

	pFrom.PendingChange -= amountMoving
	pTo.PendingChange += amountMoving
}

// func buildChainSystem(n int) *System {

// 	nodes := make([]Component, 0, n)
// 	pipes := make([]*Pipe, 0, n-1)

// 	// create n reservoirs
// 	for i := range n {
// 		id := fmt.Sprintf("R%04d", i)
// 		// each reservoir: larger area so pipes behave as small volumes
// 		cap := 1000.0
// 		area := 5.0
// 		vol := rand.Float64() * cap // random start volume
// 		// Convert Volume to Quantity (Water)
// 		quantity := vol * Water.Density

// 		res := &Reservoir{
// 			Basics: Basics{Identifier: id[0:1], Color: [3]byte{byte(rand.Intn(256)), byte(rand.Intn(256)), byte(rand.Intn(256))}},
// 			Structurals: Structurals{
// 				MaxVolume: cap,
// 				Area:      area,
// 				Quantity:  quantity,
// 				Contents:  []MaterialDef{Water},
// 			},
// 		}
// 		nodes = append(nodes, res)
// 	}

// 	// connect them in a simple chain: node[i] -> node[i+1]
// 	for i := 0; i < n-1; i++ {
// 		// diameter and length tuned for performance; adjust as needed
// 		radius := 1.0 + rand.Float64()
// 		length := 0.5 + rand.Float64()*2.0
// 		p := NewPipe(nodes[i], nodes[i+1], radius, length)
// 		pipes = append(pipes, p)
// 	}

// 	sys := &System{
// 		Nodes: nodes,
// 		Pipes: pipes,
// 		Ticks: 0,
// 	}
// 	// initialise total head/pressure for all nodes
// 	for _, n := range sys.Nodes {
// 		if s := n.GetStructurals(); s != nil {
// 			if s.Area > 0 {
// 				// s.Pressure = s.Volume / s.Area // Deprecated, calculated dynamically
// 				s.Pressure = TotalHead(n) // just for debug/init
// 			}
// 		}
// 	}
// 	return sys
// }
