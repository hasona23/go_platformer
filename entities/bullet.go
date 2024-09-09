package entities

import (
	"go_platformer/assets"

	"go_platformer/spark"
)

type Bullet struct {
	*spark.PhysicsEntity
	Dead bool
}

func NewBullet(x, y float32, speed int) *Bullet {
	bullet := &Bullet{}
	bullet.PhysicsEntity = spark.NewPhyscisEntity(x, y, 1, float32(speed), "bullet")
	bullet.Vel.Dir[0] = 1
	bullet.Sprite = *spark.NewAnimeSprite(assets.SpriteSheet, 16, 16)
	bullet.Sprite.Add(spark.NewAnimationFrame(4, 5, 4, 5, 10, "main"))
	bullet.Dead = false

	return bullet
}
func (b *Bullet) Update(tiles map[[2]int]spark.Rect, enemies []*Enemy) {
	b.Move(tiles)
	for _, i := range b.Collisions {
		if i {
			b.Dead = true
		}
	}
	for _, e := range enemies {
		if e.Collider().Collide(b.Collider()) {
			e.Dead = true
			b.Dead = true
		}
	}
}
