package main

import (
	"fmt"
	"log"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	screenWidth  = 550
	screenHeight = 320
)

func main() {
	const numTanks = 2000000
	// Assuming buildChainSystem and System field are removed,
	// the Game struct initialization needs to be adjusted.
	// For now, initializing g without System.
	g := &Game{System: buildChainSystem(numTanks)}
	fmt.Printf("--- Starting Simulation with %d tanks ---\n", numTanks)

	// --- Run Game ---
	ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
	ebiten.SetWindowTitle("My generator game")
	time.Sleep(5 * time.Second)
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
