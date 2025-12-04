package main

import (
	"fmt"
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

	// op := &text.DrawOptions{}
	// op.LayoutOptions.LineSpacing = 15
	// op.GeoM.Translate(0, 60)
	// op.ColorScale.ScaleWithColor(color.White)
	// text.Draw(screen, g.System.Status, normalFont, op)

	debugOptions := &text.DrawOptions{}
	// op.GeoM.Translate(0, 5)
	// op.ColorScale.ScaleWithColor(color.Gray{})
	text.Draw(screen, fmt.Sprintf("Tick %d | TPS %.2f | FPS %.2f\n", g.System.Ticks, ebiten.ActualTPS(), ebiten.ActualFPS()), normalFont, debugOptions)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func buildChainSystem(n int) *simulation.System {

	nodes := make([]components.Component, 0, n)
	pipes := make([]*components.Pipe, 0, n-1)

	// create n reservoirs
	for i := range n {
		id := fmt.Sprintf("R%04d", i)
		// each reservoir: larger area so pipes behave as small volumes
		cap := 1000.0
		area := 5.0
		vol := rand.Float64() * cap // random start volume
		res := &components.Reservoir{
			Basics: components.Basics{Identifier: id[0:1], Color: [3]byte{byte(rand.Intn(256)), byte(rand.Intn(256)), byte(rand.Intn(256))}},
			Structurals: components.Structurals{
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
		p := components.NewPipe(nodes[i], nodes[i+1], radius, length)
		pipes = append(pipes, p)
	}

	sys := &simulation.System{
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

func main() {
	const numTanks = 2000000
	sys := buildChainSystem(numTanks)

	g := &Game{
		System: sys,
	}

	fmt.Printf("--- Starting Simulation with %d tanks ---\n", numTanks)

	// --- Run Game ---
	ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
	ebiten.SetWindowTitle("My generator game")
	time.Sleep(5 * time.Second)
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
