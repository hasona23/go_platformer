package entities

import (
	"go_platformer/assets"

	"math"
	"math/rand"

	"go_platformer/spark"
)

type Enemy struct {
	*spark.PhysicsEntity
	Dead bool
}

func NewEnemy(x, y float32) *Enemy {
	enemy := &Enemy{}
	enemy.PhysicsEntity = spark.NewPhyscisEntity(x, y, 1, 1, "enemy")
	enemy.Sprite = *spark.NewAnimeSprite(assets.SpriteSheet, 16, 16)
	e1 := (spark.NewAnimationFrame(1, 3, 5, 6, .5, "yellowbee"))
	e2 := (spark.NewAnimationFrame(3, 5, 5, 6, .5, "bluebee"))
	e3 := (spark.NewAnimationFrame(5, 7, 5, 6, .5, "worm"))
	switch rand.Int31n(4) {
	case 1:
		enemy.Sprite.Add(e1)
		enemy.Vel.Dir = [2]float32{1, 0}
		enemy.Pos[1] -= 4
	case 2:
		enemy.Sprite.Add(e2)
		enemy.Vel.Dir = [2]float32{1, 0}
		enemy.Pos[1] -= 4
	default:
		enemy.Sprite.Add(e3)
		enemy.Vel.Dir = [2]float32{1, 1}

	}
	enemy.Dead = false
	return enemy
}
func (e *Enemy) Update(tiles map[[2]int]spark.Rect, player *Player) {

	e.PhysicsEntity.Move(tiles)

	_, ok := e.GetAroundTilesMap(tiles)[[2]int{int(math.Ceil(float64(e.Collider().X)/float64(e.Sprite.GetWidth())) + float64(e.Vel.Dir[0])),
		int(math.Ceil(float64(e.Collider().Y)/float64(e.Sprite.GetHeight()) + 1))}]

	if e.Collisions[spark.Right] || e.Collisions[spark.Left] || !ok {
		e.Vel.Dir[0] *= -1
		if e.Sprite.Effect == spark.None {
			e.Sprite.Effect = spark.FlipHorizontal
		} else {
			e.Sprite.Effect = spark.None
		}
	}
	if e.Collider().Collide(player.Collider()) {
		if player.Vel.Dir[1] > 0 && player.Collider().Y < e.Collider().Y {
			e.Dead = true
		} else {
			player.Died = true
		}
	}

}
