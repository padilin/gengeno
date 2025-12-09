package main

import "github.com/hajimehoshi/ebiten/v2"

type Entity struct {
	X, Y      int
	Component Component
	Sprite    *Sprite
}

func NewEntity(x, y int, comp Component, img *ebiten.Image, drawOrder int) *Entity {
	return &Entity{
		X:         x,
		Y:         y,
		Component: comp,
		Sprite: &Sprite{
			Image:     img,
			DrawOrder: drawOrder,
		},
	}
}
