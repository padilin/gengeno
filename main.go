package main

import (
	"fmt"
	"image/color"
	"log"
	"math/rand"
	"time"

	"gengeno/components"
	"gengeno/simulation"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"golang.org/x/image/font/basicfont"
)

// Create a text.Face that can be reused.
// text.NewStdFace is in the text/v2 package.
var normalFont text.Face = text.NewGoXFace(basicfont.Face7x13)

type Factory struct {
	tiles  []bool
	width  int
	height int
}

func NewFactory(width, height int, overall int) *Factory {
	f := &Factory{
		tiles:  make([]bool, width*height),
		width:  width,
		height: height,
	}
	f.init(overall)
	return f
}

func (f *Factory) init(overallSize int) {
	for range overallSize {
		x := rand.Intn(f.width)
		y := rand.Intn(f.height)
		f.tiles[y*f.width+x] = true
	}
}

// One tick update
func (f *Factory) Update() {
	width := f.width
	height := f.height
	next := make([]bool, width*height)
	// Stuff
	f.tiles = next
}

func (f *Factory) Draw(pix []byte) {

}

const (
	screenWidth  = 550
	screenHeight = 320
)

// Entity now holds a component and its position.
type Entity struct {
	x         int
	y         int
	component components.Component
}

type Game struct {
	System *simulation.System
}

func (g *Game) Update() error {
	g.System.Tick()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// if g.pixels == nil {
	// 	g.pixels = make([]byte, screenWidth*screenHeight*4)
	// }
	// g.factory.Draw(g.pixels)
	// screen.WritePixels(g.pixels)

	// // Draw each entity on the board
	// const tileSize = 16 // Let's define a size for our grid cells
	// for _, e := range g.board {
	// 	drawX := e.x * tileSize
	// 	// The Y position for drawing text is the baseline, so we add the tile size.
	// 	drawY := e.y*tileSize + tileSize

	// 	identifier := e.component.GetIdentifier()
	// 	r, g, b := e.component.GetColor()

	// 	options := &text.DrawOptions{}
	// 	options.GeoM.Translate(float64(drawX), float64(drawY))

	// 	// Set color by scaling the ColorM.
	// 	// The values are normalized from 0-255 to 0.0-1.0.
	// 	options.ColorScale = ebiten.ColorScale{}
	// 	options.ColorScale.ScaleWithColor(color.RGBA{r, g, b, 255})

	// 	// Draw the identifier string at the entity's position with its color.
	// 	text.Draw(screen, identifier, normalFont, options)
	// }

	op := &text.DrawOptions{}
	op.LayoutOptions.LineSpacing = 15
	op.GeoM.Translate(0, 60)
	op.ColorScale.ScaleWithColor(color.White)
	text.Draw(screen, g.System.Status, normalFont, op)

	// --- Debug Text for Reservoir B ---
	// Find Reservoir B to display its pressure for debugging. This is safer than a long chain of assertions.
	//var resBPressure float64
	// for _, e := range g.board {
	// 	if e.component.GetIdentifier() == "B" {
	// 		if s := e.component.GetStructurals(); s != nil {
	// 			resBPressure = s.Pressure
	// 		}
	// 		break
	// 	}
	// }

	// op2 := &text.DrawOptions{}
	// op2.GeoM.Translate(0, 120) // Use op2 here, and position it below the first status text.
	// op2.ColorScale.ScaleWithColor(color.White)
	// some_text := fmt.Sprintf("Tick %d |", g.System.Ticks)
	// text.Draw(screen, some_text, normalFont, op2)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	// // --- Component Setup ---
	// resA := &components.Reservoir{
	// 	Basics: components.Basics{Identifier: "A", Color: [3]byte{0, 0, 255}},
	// 	Structurals: components.Structurals{
	// 		MaxCapacity:     1000,
	// 		CurrentCapacity: 150,
	// 		Area:            20.0,
	// 	},
	// }
	// resB := &components.Reservoir{
	// 	Basics: components.Basics{Identifier: "B", Color: [3]byte{0, 0, 255}},
	// 	Structurals: components.Structurals{
	// 		MaxCapacity:     2000,
	// 		CurrentCapacity: 0,
	// 		Area:            20.0,
	// 	},
	// }
	// pipe1 := &components.Pipe{
	// 	Basics:      components.Basics{Identifier: "=", Color: [3]byte{128, 128, 128}},
	// 	FlowArea:    1.0,
	// 	Structurals: components.Structurals{MaxCapacity: 10.0, CurrentCapacity: 0.0, Area: 1.0},
	// }
	// pipe2 := &components.Pipe{
	// 	Basics:      components.Basics{Identifier: "=", Color: [3]byte{128, 128, 128}},
	// 	FlowArea:    1.0,
	// 	Structurals: components.Structurals{MaxCapacity: 10.0, CurrentCapacity: 0.0, Area: 1.0},
	// }
	// resA.Outputs = append(resA.Outputs, pipe1)
	// pipe1.Inputs = append(pipe1.Inputs, resA)
	// pipe1.Outputs = append(pipe1.Outputs, pipe2)
	// pipe2.Inputs = append(pipe2.Inputs, pipe1)
	// pipe2.Outputs = append(pipe2.Outputs, resB)
	// resB.Inputs = append(resB.Inputs, pipe2)

	// // --- Game and Simulation Initialization ---
	// // The board holds all entities that need to be drawn.
	// board := []*Entity{
	// 	{x: 5, y: 5, component: resA},
	// 	{x: 6, y: 5, component: pipe1},
	// 	{x: 7, y: 5, component: pipe2},
	// 	{x: 8, y: 5, component: resB},
	// }

	// // The simulation only needs the components that are active.
	// // In a larger game, you might only add pipes, pumps, etc., to this list.
	// allComponentsOnBoard := []components.Component{resA, pipe1, pipe2, resB}

	// connections := []simulation.Connection{
	// 	{A: resA, B: pipe1, FlowArea: pipe1.FlowArea},
	// 	{A: pipe1, B: pipe2, FlowArea: (pipe1.FlowArea+pipe2.FlowArea)/2},
	// 	{A: pipe2, B: resB, FlowArea: pipe2.FlowArea},
	// }
	// simulation := simulation.NewSimulation(connections, allComponentsOnBoard)

	// for _, e := range board {
	// 	s_structs := e.component.GetStructurals()
	// 	if s_structs.Area > 0 {
	// 		s_structs.Pressure = s_structs.CurrentCapacity / s_structs.Area
	// 	}
	// }

	// g := &Game{
	// 	factory:    NewFactory(screenWidth, screenHeight, int((screenWidth*screenHeight)/10)),
	// 	board:      board,
	// 	simulation: simulation,
	// }

	// --- The Setup ---
	// 1. Source Tank (Large, Full)
	source := &components.Reservoir{
		Basics: components.Basics{Identifier: "A", Color: [3]byte{0, 0, 255}},
		Structurals: components.Structurals{
			MaxCapacity:     1000,
			CurrentCapacity: 150,
			Area:            5.0,
			Volume:          50.0,
		},
	}

	// 2. The T-Junction
	// IMPORTANT: Give it a small area (like a pipe standing up).
	// If Area is too small, simulation oscillates. If too big, it's laggy.
	// 0.05 is roughly a 25cm pipe.
	junction := &components.Reservoir{
		Basics: components.Basics{Identifier: "T", Color: [3]byte{0, 0, 255}},
		Structurals: components.Structurals{
			MaxCapacity:     1000,
			CurrentCapacity: 150,
			Area:            0.05,
			Volume:          0.0,
			IsJunction:      true,
		},
	}

	// 3. Two Destination Tanks (Empty)
	dest1 := &components.Reservoir{
		Basics: components.Basics{Identifier: "D1", Color: [3]byte{0, 0, 255}},
		Structurals: components.Structurals{
			MaxCapacity:     1000,
			CurrentCapacity: 150,
			Area:            5.0,
			Volume:          0.0,
		},
	}
	dest2 := &components.Reservoir{
		Basics: components.Basics{Identifier: "D2", Color: [3]byte{0, 0, 255}},
		Structurals: components.Structurals{
			MaxCapacity:     1000,
			CurrentCapacity: 150,
			Area:            5.0,
			Volume:          0.0,
		},
	}

	// --- The Pipes ---
	// Pipe 1: Source -> Junction (Feeder)
	p1 := components.NewPipe(source, junction, 10.0, 0.2)

	// Pipe 2: Junction -> Dest1 (Branch A)
	p2 := components.NewPipe(junction, dest1, 10.0, 0.2)

	// Pipe 3: Junction -> Dest2 (Branch B)
	// Let's make this pipe longer, so it should fill slower!
	p3 := components.NewPipe(junction, dest2, 20.0, 0.2)

	sys := simulation.System{
		Nodes:  []components.Component{source, junction, dest1, dest2},
		Pipes:  []*components.Pipe{p1, p2, p3},
		Ticks:  0,
		Status: "Ready",
	}

	g := &Game{
		System: &sys,
	}

	fmt.Println("--- Starting T-Junction Simulation ---")

	// --- Run Game ---
	ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
	ebiten.SetWindowTitle("My generator game")
	time.Sleep(5 * time.Second)
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
