package game

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

func (t *Tile) Entities() []*Entity {
	return t.entities
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
