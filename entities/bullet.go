package entities

import (
	"go_platformer/assets"
	"go_platformer/components"
)

type Bullet struct {
	*PhysicsEntity
	Dead bool
}

func NewBullet(x, y, speed int) *Bullet {
	bullet := &Bullet{}
	bullet.PhysicsEntity = NewPhyscisEntity(x, y, 1, float32(speed), "bullet")
	bullet.Vel.Dir[0] = 1
	bullet.Anim = components.NewAnimeSprite(assets.SpriteSheet, 16, 16)
	bullet.Anim.Add(components.NewAnimationFrame(4, 5, 4, 5, 2, "main"))
	bullet.Dead = false

	return bullet
}
func (b *Bullet) Update(tiles map[[2]int]components.Rect, enemies []*Enemy) {
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
