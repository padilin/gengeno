package simulation

import (
	"fmt"
	"math"

	"gengeno/components"
)

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
