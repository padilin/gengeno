package game

import (
	"fmt"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	System               *System
	w, h                 int
	currentLevel         *Level
	camX, camY           float64
	camScale             float64
	camScaleTo           float64
	mousePanX, mousePanY int
	offscreen            *ebiten.Image
	pause                bool
}

func (g *Game) SetPause(p bool) {
	g.pause = p
}

func NewGame() (*Game, error) {
	// TODO: Move system initialization here
	_, err := LoadSpriteSheet(32)
	if err != nil {
		log.Fatal(err)
	}

	g := &Game{
		currentLevel: nil,
		camScale:     1,
		camScaleTo:   1,
		mousePanX:    0,
		mousePanY:    0,
		pause:        true,
	}
	l, err := NewLevel(g)
	if err != nil {
		return nil, fmt.Errorf("failed to create new level: %s", err)
	}
	g.currentLevel = l
	return g, nil
}

func (g *Game) Update() error {
	if !g.pause {
		g.System.Tick()
	}
	if ebiten.IsKeyPressed(ebiten.KeyP) {
		g.pause = !g.pause
	}

	// Target scroll zoom level.
	var scrollY float64
	if ebiten.IsKeyPressed(ebiten.KeyC) || ebiten.IsKeyPressed(ebiten.KeyPageDown) {
		scrollY = -0.25
	} else if ebiten.IsKeyPressed(ebiten.KeyE) || ebiten.IsKeyPressed(ebiten.KeyPageUp) {
		scrollY = 0.25
	} else {
		_, scrollY = ebiten.Wheel()
		if scrollY < -1 {
			scrollY = -1
		} else if scrollY > 1 {
			scrollY = 1
		}
	}
	g.camScaleTo += scrollY

	// Clamp scale.
	if g.camScaleTo < 0.01 {
		g.camScaleTo = 0.01
	} else if g.camScaleTo > 100 {
		g.camScaleTo = 100
	}

	// Smooth screen transitions.
	div := 10.0
	if g.camScaleTo > g.camScale {
		g.camScale += (g.camScaleTo - g.camScale) / div
	} else {
		g.camScale -= (g.camScale - g.camScaleTo) / div
	}

	// Pan caameria via keyboard.
	pan := 7.0 / g.camScale
	if ebiten.IsKeyPressed(ebiten.KeyA) || ebiten.IsKeyPressed(ebiten.KeyLeft) {
		g.camX -= pan
	} else if ebiten.IsKeyPressed(ebiten.KeyD) || ebiten.IsKeyPressed(ebiten.KeyRight) {
		g.camX += pan
	} else if ebiten.IsKeyPressed(ebiten.KeyW) || ebiten.IsKeyPressed(ebiten.KeyUp) {
		g.camY += pan
	} else if ebiten.IsKeyPressed(ebiten.KeyS) || ebiten.IsKeyPressed(ebiten.KeyDown) {
		g.camY -= pan
	}

	// Pan camera via mouse.
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight) {
		if g.mousePanX == math.MinInt32 && g.mousePanY == math.MinInt32 {
			g.mousePanX, g.mousePanY = ebiten.CursorPosition()
		} else {
			x, y := ebiten.CursorPosition()
			dx, dy := float64(g.mousePanX-x)*(pan*100), float64(g.mousePanY-y)*(pan*100)
			g.camX, g.camY = g.camX-dx, g.camY+dy
		}
	} else if g.mousePanX != math.MinInt32 && g.mousePanY != math.MinInt32 {
		g.mousePanX, g.mousePanY = math.MinInt32, math.MinInt32
	}

	// Clamp camera position
	worldWidth := float64(g.currentLevel.Width * g.currentLevel.tileSize / 2)
	worldHeight := float64(g.currentLevel.Height * g.currentLevel.tileSize / 2)
	if g.camX < -worldWidth {
		g.camX = -worldWidth
	} else if g.camX > worldWidth {
		g.camX = worldWidth
	}
	if g.camY < -worldHeight {
		g.camY = -worldHeight
	} else if g.camY > worldHeight {
		g.camY = worldHeight
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.renderLevel(screen)
	// ebitenutil.DebugPrint(screen, fmt.Sprintf("Fill: %.1f", g.System.Nodes[0].GetStructurals().CurrentCapacity))

	// ebitenutil.DebugPrint(screen, fmt.Sprintf("KEYS WASD EC R\nFPS  %0.0f\nTPS  %0.0f\nSCA  %0.2f\nPOS  %0.0f,%0.0f", ebiten.ActualFPS(), ebiten.ActualTPS(), g.camScale, g.camX, g.camY))
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	g.w, g.h = outsideWidth, outsideHeight
	return g.w, g.h
}

func (g *Game) CartesianToIso(x, y float64) (float64, float64) {
	tileSize := g.currentLevel.tileSize
	ix := (x - y) * float64(tileSize/2)
	iy := (x + y) * float64(tileSize/4)
	return ix, iy
}

func (g *Game) IsoToCartesian(x, y float64) (float64, float64) {
	tileSize := g.currentLevel.tileSize
	cx := (x/float64(tileSize/2) + y/float64(tileSize/4)) / 2
	cy := (y/float64(tileSize/4) - x/float64(tileSize/2)) / 2
	return cx, cy
}

func (g *Game) renderLevel(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	padding := float64(g.currentLevel.tileSize) * g.camScale
	cx, cy := float64(g.w/2), float64(g.h/2)

	scaleLater := g.camScale > 1
	target := screen
	scale := g.camScale

	if scaleLater {
		if g.offscreen != nil {
			if g.offscreen.Bounds().Size() != screen.Bounds().Size() {
				g.offscreen.Deallocate()
				g.offscreen = nil
			}
		}
		if g.offscreen == nil {
			s := screen.Bounds().Size()
			g.offscreen = ebiten.NewImage(s.X, s.Y)
		}
		target = g.offscreen
		target.Clear()
		scale = 1
	}

	for y := 0; y < g.currentLevel.Height; y++ {
		for x := 0; x < g.currentLevel.Width; x++ {
			xi, yi := g.CartesianToIso(float64(x), float64(y))

			// Skip offscreen
			drawX, drawY := ((xi-g.camX)*g.camScale)+cx, ((yi+g.camY)*g.camScale)+cy
			if drawX+padding < 0 || drawY+padding < 0 || drawX > float64(g.w) || drawY > float64(g.h) {
				continue
			}

			t := g.currentLevel.tiles[y][x]
			if t == nil {
				continue
			}

			op.GeoM.Reset()
			op.GeoM.Translate(xi, yi)
			op.GeoM.Translate(-g.camX, g.camY)
			op.GeoM.Scale(scale, scale)
			op.GeoM.Translate(cx, cy)

			t.Draw(target, op)
		}
	}

	if scaleLater {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(-cx, -cy)
		op.GeoM.Scale(float64(g.camScale), float64(g.camScale))
		op.GeoM.Translate(cx, cy)
		screen.DrawImage(target, op)
	}
}
