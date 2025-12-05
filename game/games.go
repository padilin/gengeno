package main

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type Game struct {
	System *System
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
	debugOptions.ColorScale.ScaleWithColor(color.White) // Ensure text is visible
	// op.GeoM.Translate(0, 5)
	// op.ColorScale.ScaleWithColor(color.Gray{})
	text.Draw(screen, fmt.Sprintf("Tick %d | TPS %.2f | FPS %.2f\n", g.System.Ticks, ebiten.ActualTPS(), ebiten.ActualFPS()), normalFont, debugOptions)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}
