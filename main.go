package main

import (
	"log"
	"math/rand"

	"gengeno/components"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"golang.org/x/image/font/basicfont"
	"image/color"
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
	screenWidth  = 320
	screenHeight = 240
)

// Entity now holds a component and its position.
type Entity struct {
	x         int
	y         int
	component components.Component
}

type Game struct {
	factory *Factory
	pixels  []byte
	board   []*Entity
}

func (g *Game) Update() error {
	g.factory.Update()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	if g.pixels == nil {
		g.pixels = make([]byte, screenWidth*screenHeight*4)
	}
	g.factory.Draw(g.pixels)
	screen.WritePixels(g.pixels)

	// Draw each entity on the board
	const tileSize = 16 // Let's define a size for our grid cells
	for _, e := range g.board {
		drawX := e.x * tileSize
		// The Y position for drawing text is the baseline, so we add the tile size.
		drawY := e.y*tileSize + tileSize

		identifier := e.component.GetIdentifier()
		r, g, b := e.component.GetColor()

		options := &text.DrawOptions{}
		options.GeoM.Translate(float64(drawX), float64(drawY))

		// Set color by scaling the ColorM.
		// The values are normalized from 0-255 to 0.0-1.0.
		options.ColorScale = ebiten.ColorScale{}
		options.ColorScale.ScaleWithColor(color.RGBA{r, g, b, 255})

		// Draw the identifier string at the entity's position with its color.
		text.Draw(screen, identifier, normalFont, options)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	// Create an instance of a component, e.g., a Reservoir.
	// This is where you can define the properties of your component.
	componentA := &components.Reservoir{
		Basics: components.Basics{
			Identifier: "A",
			Color:      [3]byte{255, 0, 0}, // Red
		},
	}

	g := &Game{
		factory: NewFactory(screenWidth, screenHeight, int((screenWidth*screenHeight)/10)),
		board: []*Entity{
			// This entity is now linked to your componentA instance.
			{x: 5, y: 5, component: componentA},
		},
	}

	ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
	ebiten.SetWindowTitle("My generator game")
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
