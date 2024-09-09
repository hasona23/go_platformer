package entities

import (
	"go_platformer/assets"
	"go_platformer/spark/particles"
	"image"
	"image/color"
	"slices"

	"go_platformer/spark"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Player struct {
	*spark.PhysicsEntity
	Bullets       []*Bullet
	shootingTimer spark.Timer
	Ammo          int
	Died          bool
	isJumping     bool
	isShooting    bool
	//	shootParticles  *particles.ParticleSystem
	bulletParticles *particles.ParticleSystem
}

func NewPlayer() *Player {
	player := &Player{}
	player.PhysicsEntity = spark.NewPhyscisEntity(40, 330, 1, 1, "player")

	player.Sprite = *spark.NewAnimeSprite(assets.SpriteSheet, 16, 16)
	player.Sprite.Add(spark.NewAnimationFrame(0, 1, 4, 5, 2, "idle"))
	player.Sprite.Add(spark.NewAnimationFrame(0, 2, 4, 5, .2, "run"))
	player.Sprite.Add(spark.NewAnimationFrame(2, 3, 4, 5, 0.7, "shoot"))
	player.isJumping = false
	player.isShooting = false
	player.Ammo = 10
	player.shootingTimer = spark.NewTimer(.3)
	color := color.RGBA{255, 186, 57, 255}
	player.bulletParticles = particles.NewParticleSystem(particles.WithArea(spark.NewRect(int(player.Pos[0]), int(player.Pos[1]), 20, 20)),
		particles.WithMotionType(particles.Outward), particles.WithShrinking(0.1),
		particles.WithModelParticle(*particles.NewParticle(particles.WithColor(color), particles.WithScale(4), particles.WithSpeed(1))))

	return player
}
func Lerp(a, b, t float32) float32 {
	return a + (b-a)*t
}

func (p *Player) Update(tiles map[[2]int]spark.Rect) {
	if p.Collisions[spark.Down] {
		p.isJumping = false
	}
	p.Vel.Dir[0] = 0
	//applies more gravity when falling
	if p.Vel.Dir[1] > .5 {
		p.Vel.Dir[1] = min(3, p.Vel.Dir[1]+.8)
	} else {
		p.Vel.Dir[1] = min(2, p.Vel.Dir[1]+0.1)
	}
	//Movement Controls

	if (inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) || inpututil.IsKeyJustPressed(ebiten.KeyW) || inpututil.IsKeyJustPressed(ebiten.KeySpace)) && !p.isJumping {
		p.Vel.Dir[1] = -3
		p.isJumping = true
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) || ebiten.IsKeyPressed(ebiten.KeyA) {
		p.Vel.Dir[0] = -1
		p.Sprite.Effect = spark.FlipHorizontal
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) || ebiten.IsKeyPressed(ebiten.KeyD) {
		p.Vel.Dir[0] = 1
		p.Sprite.Effect = spark.None
	}
	//Change animation based on motion
	if !p.isShooting {
		if p.Vel.Dir[0] != 0 {
			p.Sprite.ChangeAnim("run")
		} else if p.Vel.Dir[0] == 0 {
			p.Sprite.ChangeAnim("idle")
		}
	}
	p.PhysicsEntity.Move(tiles)
	p.shoot()

	//p.shootParticles.Update()
	//p.bulletParticles.Area.X = int(p.Pos[0])
	//p.bulletParticles.Area.Y = int(p.Pos[1])
	p.bulletParticles.Update()
	p.Bullets = slices.DeleteFunc(p.Bullets, func(b *Bullet) bool { return b.Dead })

}
func (p Player) Draw(screen *ebiten.Image, cam spark.Cam) {
	for _, b := range p.Bullets {
		b.Draw(screen, cam)
	}
	//draws effect of gun when shooting
	if p.isShooting {
		op := &ebiten.DrawImageOptions{}
		if p.Sprite.Effect == spark.None {
			op.GeoM.Translate(float64(p.Pos[0]+cam.X+16), float64(float64(p.Pos[1]+cam.Y)+2))
		} else if p.Sprite.Effect == spark.FlipHorizontal {
			p.Sprite.FlipHorizontal(op)
			op.GeoM.Translate(float64(p.Pos[0]+cam.X-16), float64(float64(p.Pos[1]+cam.Y))+2)

		}
		screen.DrawImage(assets.SpriteSheet.SubImage(image.Rect(3*16, 4*16, 4*16, 5*16)).(*ebiten.Image), op)
	}
	p.PhysicsEntity.Draw(screen, cam)
	//p.shootParticles.DrawCam(screen, cam)
	//p.bulletParticles.Draw(screen)
	p.bulletParticles.DrawCam(screen, cam)

}
func (p *Player) shoot() {
	if p.Ammo <= 0 {
		p.isShooting = false
		return
	}
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) || inpututil.IsKeyJustPressed(ebiten.KeyE) {
		p.isShooting = true
		p.Sprite.ChangeAnim("shoot")
		bSpeed := 0
		//	p.shootParticles.Area.Y = int(p.Pos.Pos[1]) + 2
		if p.Sprite.Effect == spark.FlipHorizontal {
			bSpeed = -5
			//		p.shootParticles.Area.X = p.Pos.Pos[0] - 16
		} else {
			//		p.shootParticles.Area.X = p.Pos.Pos[0] + 16
			bSpeed = 5
		}
		//	p.shootParticles.Spawn(10)
		p.Bullets = append(p.Bullets, NewBullet(p.Pos[0], p.Pos[1], bSpeed))
		p.Ammo--
	}
	if p.shootingTimer.Ticked() && p.isShooting {
		p.isShooting = false
	}
	if p.isShooting {
		p.shootingTimer.UpdateTimerTPS()
	}

}
func (p *Player) UpdateBullets(tiles map[[2]int]spark.Rect, enemies []*Enemy) {
	for i, b := range p.Bullets {
		p.Bullets[i].Update(tiles, enemies)
		if b.Dead {
			p.bulletParticles.Area.X = b.Collider().X + b.Collider().Width
			p.bulletParticles.Area.Y = b.Collider().Y + b.Collider().Height/2
			p.bulletParticles.Spawn(30)
		}
	}
}
