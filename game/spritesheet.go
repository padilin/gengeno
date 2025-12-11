// Copyright 2021 The Ebiten Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"bytes"
	"image"
	_ "image/png"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
)

var spriteSet map[string]*Sprite

// SpriteSheet represents a collection of sprite images.
type SpriteSheet struct {
	Floor           *ebiten.Image
	Weird           *ebiten.Image
	Reservoir1      *ebiten.Image
	Reservoir2      *ebiten.Image
	PipeDown        *ebiten.Image
	PipeHorz        *ebiten.Image
	PipeEnterLeft   *ebiten.Image
	PipeEnterRight  *ebiten.Image
	PipeVert        *ebiten.Image
	PipeLeftToDown  *ebiten.Image
	PipeRightToDown *ebiten.Image
}

// LoadSpriteSheet loads the embedded SpriteSheet.
func LoadSpriteSheet(tileSize int) (*SpriteSheet, error) {
	data, err := os.ReadFile("assets/floor-tile-1.png")
	if err != nil {
		return nil, err
	}
	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	sheet := ebiten.NewImageFromImage(img)

	// spriteAt returns a sprite at the provided coordinates.
	spriteAt := func(x, y int) *ebiten.Image {
		return sheet.SubImage(image.Rect(x*tileSize, (y+1)*tileSize, (x+1)*tileSize, y*tileSize)).(*ebiten.Image)
	}

	// Populate SpriteSheet.
	s := &SpriteSheet{}
	// s.Floor = spriteAt(0, 0)
	// s.Weird = spriteAt(1, 0)
	// s.Reservoir1 = spriteAt(2, 0)
	// s.Reservoir2 = spriteAt(2, 1)
	// s.PipeDown = spriteAt(3, 0)
	// s.PipeHorz = spriteAt(3, 1)
	// s.PipeEnterLeft = spriteAt(3, 2)
	// s.PipeEnterRight = spriteAt(3, 3)
	// s.PipeVert = spriteAt(3, 4)
	// s.PipeLeftToDown = spriteAt(3, 5)
	// s.PipeRightToDown = spriteAt(3, 6)
	// s.Wall = spriteAt(2, 3)
	// s.Statue = spriteAt(5, 4)
	// s.Tube = spriteAt(3, 4)
	// s.Crown = spriteAt(8, 6)
	// s.Portal = spriteAt(5, 6)

	spriteSet = map[string]*Sprite{
		"floor":           {Image: spriteAt(0, 0), DrawOrder: 0},
		"pipe_enter_left": {Image: spriteAt(3, 2), DrawOrder: 10},
		"reservoir_full":  {Image: spriteAt(2, 0), DrawOrder: 5},
		"reservoir_high":  {Image: spriteAt(2, 1), DrawOrder: 5},
		"reservoir_mid":   {Image: spriteAt(2, 2), DrawOrder: 5},
		"reservoir_low":   {Image: spriteAt(2, 3), DrawOrder: 5},
		"reservoir_empty": {Image: spriteAt(2, 4), DrawOrder: 5},
	}

	return s, nil
}
