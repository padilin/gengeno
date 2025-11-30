package main

import (
	"log"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	
	
)

type Factory struct {
	tiles []bool
	width int
	height int
}

func NewFactory(width, height int, overall int) *Factory {
	f := &Factory{
		tiles: make([]bool, width*height),
		width: width,
		height: height,
	}
	f.init(overall)
	return f
}

func (f *Factory) init(overallSize int) {
	for i := 0; i < overallSize; i++ {
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
	for i, v := range f.tiles {
		if v {
			pix[4*i] = 0xff
		} else {
			pix[4*i] = 0
		}
	}
}

const (
	screenWidth = 320
	screenHeight = 240
)

type Game struct {
	factory *Factory
	pixels []byte
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
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	g := &Game{
		factory: NewFactory(screenWidth, screenHeight, int((screenWidth*screenHeight)/10)),
	}

	ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
	ebiten.SetWindowTitle("My generator game")
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
