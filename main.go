package main

import (
	"log"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/padilin/gengeno/game"
)

const (
	screenWidth  = 550
	screenHeight = 320
)

func main() {
	// Create game with the level
	g, err := game.NewGame()
	if err != nil {
		log.Fatal(err)
	}
	// --- Run Game ---
	ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
	ebiten.SetWindowTitle("My generator game")
	time.Sleep(5 * time.Second)
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
