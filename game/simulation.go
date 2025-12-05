package main

import (
	"fmt"
	"math"
	"math/rand"
)

const (
	Gravity      = 9.81
	TimeStep     = 1.0 / 6.0
	FrictionFact = 0.02
	MinorLoss    = 1.5
)

func TotalHead(c Component) float64 {
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

func ApplyPending(c Component) {
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
	Nodes []Component
	Pipes []*Pipe
	Ticks int
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
		frictionLoss := FrictionFact * (p.Length / (2 * p.Radius))
		velocity := math.Sqrt((2 * Gravity * deltaH) / (1 + frictionLoss + MinorLoss))

		// Volume to move with direction
		flowVol := velocity * p.Area * TimeStep
		amountMoving := flowVol * direction

		var source, dest Component
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
}

func buildChainSystem(n int) *System {

	nodes := make([]Component, 0, n)
	pipes := make([]*Pipe, 0, n-1)

	// create n reservoirs
	for i := range n {
		id := fmt.Sprintf("R%04d", i)
		// each reservoir: larger area so pipes behave as small volumes
		cap := 1000.0
		area := 5.0
		vol := rand.Float64() * cap // random start volume
		res := &Reservoir{
			Basics: Basics{Identifier: id[0:1], Color: [3]byte{byte(rand.Intn(256)), byte(rand.Intn(256)), byte(rand.Intn(256))}},
			Structurals: Structurals{
				MaxCapacity:     cap,
				CurrentCapacity: vol,
				Area:            area,
				Volume:          vol,
			},
		}
		nodes = append(nodes, res)
	}

	// connect them in a simple chain: node[i] -> node[i+1]
	for i := 0; i < n-1; i++ {
		// diameter and length tuned for performance; adjust as needed
		radius := 1.0 + rand.Float64()
		length := 0.5 + rand.Float64()*2.0
		p := NewPipe(nodes[i], nodes[i+1], radius, length)
		pipes = append(pipes, p)
	}

	sys := &System{
		Nodes: nodes,
		Pipes: pipes,
		Ticks: 0,
	}
	// initialise total head/pressure for all nodes
	for _, n := range sys.Nodes {
		if s := n.GetStructurals(); s != nil {
			if s.Area > 0 {
				s.Pressure = s.Volume / s.Area
			}
		}
	}
	return sys
}
