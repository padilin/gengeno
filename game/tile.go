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
	"sort"

	"github.com/hajimehoshi/ebiten/v2"
)

type Sprite struct {
	Image     *ebiten.Image
	DrawOrder int
}

// Tile represents a space with an x,y coordinate within a Level. Any number of
// sprites may be added to a Tile.
type Tile struct {
	entities []*Entity
}

func (t *Tile) AddEntity(entity *Entity) {
	if t == nil || entity == nil {
		return
	}
	t.entities = append(t.entities, entity)
}

// Draw draws the Tile on the screen using the provided options.
func (t *Tile) Draw(screen *ebiten.Image, baseOptions *ebiten.DrawImageOptions) {
	if t == nil || baseOptions == nil {
		return
	}

	// Sort entities by DrawOrder (ascending)
	sorted := make([]*Entity, len(t.entities))
	copy(sorted, t.entities)
	sort.Slice(sorted, func(i, j int) bool {
		spriteI := sorted[i].CurrentSprite()
		spriteJ := sorted[j].CurrentSprite()
		if spriteI == nil || spriteJ == nil {
			return false
		}
		return spriteI.DrawOrder < spriteJ.DrawOrder
	})

	for _, e := range sorted {
		if e == nil {
			continue
		}

		// Use CurrentSprite() to get state-based sprite if a selector is set
		sprite := e.CurrentSprite()
		if sprite == nil {
			sprite = e.Sprite // fallback to default sprite
		}
		if sprite == nil || sprite.Image == nil {
			continue
		}

		opts := *baseOptions

		screen.DrawImage(sprite.Image, &opts)
	}
}
