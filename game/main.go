package main

import (
	"log"
	"time"

	// "os"
	// "runtime/pprof"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	screenWidth  = 550
	screenHeight = 320
)

func main() {
	// f, err := os.Create("cpu.prof")
	// if err != nil {
	// 	panic(err)
	// }
	// defer f.Close()
	// pprof.StartCPUProfile(f)
	// defer pprof.StopCPUProfile()
	const numTanks = 2000000
	// Assuming buildChainSystem and System field are removed,
	// the Game struct initialization needs to be adjusted.
	// For now, initializing g without System.
	// g := &Game{System: buildChainSystem(numTanks)}
	// fmt.Printf("--- Starting Simulation with %d tanks ---\n", numTanks)
	g, err := NewGame()
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
