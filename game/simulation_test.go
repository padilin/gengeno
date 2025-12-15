package main

import (
	"math"
	"testing"
)

func TestGetMaterial(t *testing.T) {
	tests := []struct {
		name string
		c    Component
		want *MaterialDef
	}{
		{
			name: "With Contents",
			c: &Reservoir{
				Structurals: Structurals{Contents: []MaterialDef{Water}},
			},
			want: &Water,
		},
		{
			name: "Empty Contents",
			c: &Reservoir{
				Structurals: Structurals{Contents: []MaterialDef{}},
			},
			want: &Water, // Defaults to Water
		},
		{
			name: "Nil Structurals",
			c:    &Reservoir{}, // Structurals zero value has empty contents
			want: &Water,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetMaterial(tt.c); *got != *tt.want { // Compare values
				t.Errorf("GetMaterial() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTotalHead(t *testing.T) {
	// Head = (Quantity / (Density * Area)) + BaseElevation
	// Water Density = 1000
	tests := []struct {
		name string
		c    Component
		want float64
	}{
		{
			name: "Simple Water Column",
			c: &Reservoir{
				Structurals: Structurals{
					Area:          1.0,
					Quantity:      1000.0, // Should result in 1m height
					BaseElevation: 0,
					Contents:      []MaterialDef{Water},
				},
			},
			want: 1.0,
		},
		{
			name: "Water Column with Elevation",
			c: &Reservoir{
				Structurals: Structurals{
					Area:          1.0,
					Quantity:      1000.0,
					BaseElevation: 10.0,
					Contents:      []MaterialDef{Water},
				},
			},
			want: 11.0,
		},
		{
			name: "Empty Area (Divide by Zero Protection)",
			c: &Reservoir{
				Structurals: Structurals{Area: 0, BaseElevation: 5},
			},
			want: 5.0,
		},
		{
			name: "Gas Pressure",
			c: &Reservoir{
				Structurals: Structurals{
					MaxVolume:     10.0,
					Quantity:      5.0,
					Contents:      []MaterialDef{{Type: TypeGas, GasConstant: 100}},
					BaseElevation: 0,
				},
			},
			// Pressure = (5 * 100) / 10 = 50
			// Head = P / (rho_water * g) = 50 / (1000 * 9.81)
			want: 50.0 / (1000.0 * 9.81),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := TotalHead(tt.c)
			if math.Abs(got-tt.want) > 0.0001 {
				t.Errorf("TotalHead() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestApplyPending(t *testing.T) {
	r := &Reservoir{
		Structurals: Structurals{
			Quantity:      100,
			PendingChange: 10,
		},
	}
	ApplyPending(r)
	if r.Structurals.Quantity != 110 {
		t.Errorf("ApplyPending() quantity = %v, want 110", r.Structurals.Quantity)
	}
	if r.Structurals.PendingChange != 0 {
		t.Errorf("ApplyPending() pending = %v, want 0", r.Structurals.PendingChange)
	}

	// Negative Check
	r.Structurals.PendingChange = -200
	ApplyPending(r)
	if r.Structurals.Quantity != 0 {
		t.Errorf("ApplyPending() quantity underflow = %v, want 0", r.Structurals.Quantity)
	}
}

func TestSystem_Tick(t *testing.T) {
	// Setup simple specific system: Source -> Pipe -> Dest
	res1 := &Reservoir{
		Structurals: Structurals{
			Area:          10,
			Quantity:      10000,
			BaseElevation: 10, // High head
			Contents:      []MaterialDef{Water},
		},
	}
	res2 := &Reservoir{
		Structurals: Structurals{
			Area:          10,
			MaxVolume:     20000,
			Quantity:      0,
			BaseElevation: 0, // Low head
			Contents:      []MaterialDef{Water},
		},
	}
	pipe := NewPipe(res1, res2, 10, 1)

	s := &System{
		Nodes: []Component{res1, res2},
		Pipes: []*Pipe{pipe},
	}

	// Run enough ticks to trigger flow (Tick % 10 == 0)
	// We need multiple logic steps for flow to propagate: Source -> Pipe -> Dest
	// Step 1 (Tick 10): Source -> Pipe (Pipe gets pending qty)
	// Step 2 (Tick 20): Pipe (now has qty) -> Dest
	for i := 0; i <= 30; i++ {
		s.Tick()
	}

	// Check if flow occurred
	if res1.Structurals.Quantity >= 10000 {
		t.Error("Tick() Source quantity did not decrease")
	}
	// With multiple steps, dest should now have received some fluid from pipe
	if res2.Structurals.Quantity <= 0 {
		t.Errorf("Tick() Dest quantity did not increase, qty=%f", res2.Structurals.Quantity)
	}
}

func Test_calculateFlow(t *testing.T) {
	// Covered implicitly by System_Tick but let's test a specific case
	from := &Reservoir{Structurals: Structurals{Area: 1, Quantity: 1000, BaseElevation: 10, Contents: []MaterialDef{Water}}}
	to := &Reservoir{Structurals: Structurals{Area: 1, MaxVolume: 1000, Quantity: 0, BaseElevation: 0, Contents: []MaterialDef{Water}}}
	// Pipe geometry will be fallback in function if not pipe

	// We need to pass valid components.
	// If neither is pipe, calculateFlow uses defaults (radius=0.5, len=1.0)
	calculateFlow(from, to, 0)

	if from.Structurals.PendingChange >= 0 {
		t.Errorf("calculateFlow failed to decrement source pending change: %v", from.Structurals.PendingChange)
	}
	if to.Structurals.PendingChange <= 0 {
		t.Errorf("calculateFlow failed to increment dest pending change: %v", to.Structurals.PendingChange)
	}
}
