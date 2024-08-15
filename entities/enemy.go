package entities

import (
	"fmt"
	"go_platformer/assets"
	"go_platformer/components"
	"math"
	"math/rand"
)

type Enemy struct {
	*PhysicsEntity
}

func NewEnemy(x, y int) *Enemy {
	enemy := &Enemy{}
	enemy.PhysicsEntity = NewPhyscisEntity(x, y, 1, 1, "enemy")
	enemy.Anim = components.NewAnimeSprite(assets.SpriteSheet, 16, 16)
	e1 := (components.NewAnimationFrame(1, 3, 5, 6, .5, "yellowbee"))
	e2 := (components.NewAnimationFrame(3, 5, 5, 6, .2, "bluebee"))
	e3 := (components.NewAnimationFrame(5, 7, 5, 6, 0.7, "worm"))
	switch rand.Int31n(4) {
	case 1:
		enemy.Anim.Add(e1)
		enemy.Vel.Dir = [2]float32{1, 0}
		enemy.Pos.Pos[1] += 4
	case 2:
		enemy.Anim.Add(e2)
		enemy.Vel.Dir = [2]float32{1, 0}
		enemy.Pos.Pos[1] += 4
	case 3:
		enemy.Anim.Add(e3)
		enemy.Vel.Dir = [2]float32{1, 1}
	}
	return enemy
}
func (e *Enemy) Update(tiles map[[2]int]components.Rect, playerRect components.Rect) {
	e.PhysicsEntity.Move(tiles)

	_, ok := e.GetAroundTilesMap(tiles)[[2]int{int(math.Ceil(float64(e.Collider().X)/float64(e.Anim.GetWidth())) + float64(e.Vel.Dir[0])),
		int(math.Ceil(float64(e.Collider().Y)/float64(e.Anim.GetHeight()) + 1))}]

	if e.Collisions["right"] || e.Collisions["left"] || !ok {
		e.Vel.Dir[0] *= -1
		e.Anim.Flip = !e.Anim.Flip
	}
	if e.Collider().Collide(playerRect) {
		fmt.Println("Collided Player")
	}
}
