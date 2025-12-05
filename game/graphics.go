package main

import (
	"math/rand"

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

// Entity now holds a component and its position.
type Entity struct {
	x         int
	y         int
	component Component
}
