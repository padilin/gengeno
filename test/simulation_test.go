package test

import (
	"math"
	"testing"

	"github.com/padilin/gengeno/game"
)

func TestGetMaterial(t *testing.T) {
	tests := []struct {
		name string
		c    game.Component
		want *game.MaterialDef
	}{
		{
			name: "With Contents",
			c: &game.Reservoir{
				Structurals: game.Structurals{Contents: []game.MaterialDef{game.Water}},
			},
			want: &game.Water,
		},
		{
			name: "Empty Contents",
			c: &game.Reservoir{
				Structurals: game.Structurals{Contents: []game.MaterialDef{}},
			},
			want: &game.Water, // Defaults to Water
		},
		{
			name: "Nil Structurals",
			c:    &game.Reservoir{}, // Structurals zero value has empty contents
			want: &game.Water,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := game.GetMaterial(tt.c); *got != *tt.want { // Compare values
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
		c    game.Component
		want float64
	}{
		{
			name: "Simple Water Column",
			c: &game.Reservoir{
				Structurals: game.Structurals{
					Area:          1.0,
					Quantity:      1000.0, // Should result in 1m height
					BaseElevation: 0,
					Contents:      []game.MaterialDef{game.Water},
				},
			},
			want: 1.0,
		},
		{
			name: "Water Column with Elevation",
			c: &game.Reservoir{
				Structurals: game.Structurals{
					Area:          1.0,
					Quantity:      1000.0,
					BaseElevation: 10.0,
					Contents:      []game.MaterialDef{game.Water},
				},
			},
			want: 11.0,
		},
		{
			name: "Empty Area (Divide by Zero Protection)",
			c: &game.Reservoir{
				Structurals: game.Structurals{Area: 0, BaseElevation: 5},
			},
			want: 5.0,
		},
		{
			name: "Gas Pressure",
			c: &game.Reservoir{
				Structurals: game.Structurals{
					MaxVolume:     10.0,
					Quantity:      5.0,
					Contents:      []game.MaterialDef{{Type: game.TypeGas, GasConstant: 100}},
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
			got := game.TotalHead(tt.c)
			if math.Abs(got-tt.want) > 0.0001 {
				t.Errorf("TotalHead() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestApplyPending(t *testing.T) {
	r := &game.Reservoir{
		Structurals: game.Structurals{
			Quantity:      100,
			PendingChange: 10,
		},
	}
	game.ApplyPending(r)
	if r.Quantity != 110 {
		t.Errorf("ApplyPending() quantity = %v, want 110", r.Quantity)
	}
	if r.PendingChange != 0 {
		t.Errorf("ApplyPending() pending = %v, want 0", r.PendingChange)
	}

	// Negative Check
	r.PendingChange = -200
	game.ApplyPending(r)
	if r.Quantity != 0 {
		t.Errorf("ApplyPending() quantity underflow = %v, want 0", r.Quantity)
	}
}

func TestSystem_Tick(t *testing.T) {
	// Setup simple specific system: Source -> Pipe -> Dest
	res1 := &game.Reservoir{
		Structurals: game.Structurals{
			Area:          10,
			Quantity:      10000,
			BaseElevation: 10, // High head
			Contents:      []game.MaterialDef{game.Water},
		},
	}
	res2 := &game.Reservoir{
		Structurals: game.Structurals{
			Area:          10,
			MaxVolume:     20000,
			Quantity:      0,
			BaseElevation: 0, // Low head
			Contents:      []game.MaterialDef{game.Water},
		},
	}
	pipe := game.NewPipe(res1, res2, 10, 1)

	s := &game.System{
		Nodes: []game.Component{res1, res2},
		Pipes: []*game.Pipe{pipe},
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
	// t.Skip("Skipping Test_calculateFlow (unexported function)")
	// Just test basic flow: R1(100) -> Pipe -> R2(0)
	// After 1 tick (CalculateFlow called manually or via Tick)

	// Create components manually to test CalculateFlow
	r1 := &game.Reservoir{Structurals: game.Structurals{MaxVolume: 100, Quantity: 100, BaseElevation: 10, Contents: []game.MaterialDef{game.Water}}}
	r2 := &game.Reservoir{Structurals: game.Structurals{MaxVolume: 100, Quantity: 0, BaseElevation: 0, Contents: []game.MaterialDef{game.Water}}}

	// CalculateFlow(From, To, PumpHead)
	// But it requires Pipe geometry usually?
	// calculateFlow implementation:
	// if pipe, ok := to.(*Pipe) ...
	// If neither is pipe, uses fallback 1.0, 1.0, 0.5

	game.CalculateFlow(r1, r2, 0)

	if r1.Structurals.PendingChange >= 0 {
		t.Error("r1 should have negative PendingChange")
	}
	if r2.Structurals.PendingChange <= 0 {
		t.Error("r2 should have positive PendingChange")
	}
}
