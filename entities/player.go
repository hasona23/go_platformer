package entities

import (
	"go_platformer/assets"
	"go_platformer/components"
	"go_platformer/particles"
	"image"
	"image/color"
	"slices"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Player struct {
	*PhysicsEntity
	Bullets       []*Bullet
	shootingTimer components.Timer
	Ammo          int
	Died          bool
	isJumping     bool
	isShooting    bool
	//	shootParticles  *particles.ParticleSystem
	bulletParticles *particles.ParticleSystem
}

func NewPlayer() *Player {
	player := &Player{}
	player.PhysicsEntity = NewPhyscisEntity(40, 330, 1, 1, "player")
	player.Anim = components.NewAnimeSprite(assets.SpriteSheet, 16, 16)
	player.Anim.Add(components.NewAnimationFrame(0, 1, 4, 5, 2, "idle"))
	player.Anim.Add(components.NewAnimationFrame(0, 2, 4, 5, .2, "run"))
	player.Anim.Add(components.NewAnimationFrame(2, 3, 4, 5, 0.7, "shoot"))
	player.isJumping = false
	player.isShooting = false
	player.Ammo = 10
	player.shootingTimer = components.NewTimer(.3)
	color := color.RGBA{255, 186, 57, 255}
	//	player.shootParticles = particles.NewParticleSystem(player.Pos.Pos[0], player.Pos.Pos[1], 4, 4, 0, particles.Outward,
	//		*particles.NewParticle(0, 0, particles.WithScale(4), particles.WithSpeed(0), particles.WithColor(color)))
	//	player.shootParticles.Decelration = 0.4
	player.bulletParticles = particles.NewParticleSystem(player.Pos.Pos[0], player.Pos.Pos[1], 4, 4, 0, particles.Outward,
		*particles.NewParticle(0, 0, particles.WithColor(color), particles.WithScale(4), particles.WithSpeed(3)))
	player.bulletParticles.Decelration = 0.2

	return player
}
func Lerp(a, b, t float32) float32 {
	return a + (b-a)*t
}

func (p *Player) Update(tiles map[[2]int]components.Rect) {
	if p.Collisions[Down] {
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
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) || ebiten.IsKeyPressed(ebiten.KeyD) {
		p.Vel.Dir[0] = 1
	}
	//Change animation based on motion
	if !p.isShooting {
		if p.Vel.Dir[0] != 0 {
			p.Anim.ChangeAnim("run")
		} else if p.Vel.Dir[0] == 0 {
			p.Anim.ChangeAnim("idle")
		}
	}
	p.PhysicsEntity.Move(tiles)
	p.shoot()
	//p.shootParticles.Update()
	p.bulletParticles.Update()
	p.Bullets = slices.DeleteFunc(p.Bullets, func(b *Bullet) bool { return b.Dead })

}
func (p Player) Draw(screen *ebiten.Image, cam components.Camera) {
	for _, b := range p.Bullets {
		b.Draw(screen, cam)
	}
	//draws effect of gun when shooting
	if p.isShooting {
		op := &ebiten.DrawImageOptions{}
		if !p.Anim.Flip {
			op.GeoM.Translate(float64(p.Pos.Pos[0]+cam.X+16), float64(float64(p.Pos.Pos[1]+cam.Y)+2))
		} else {

			p.Anim.FlipHorizontal(op)
			op.GeoM.Translate(float64(p.Pos.Pos[0]+cam.X-16), float64(float64(p.Pos.Pos[1]+cam.Y))+2)

		}
		screen.DrawImage(assets.SpriteSheet.SubImage(image.Rect(3*16, 4*16, 4*16, 5*16)).(*ebiten.Image), op)
	}
	p.PhysicsEntity.Draw(screen, cam)
	//p.shootParticles.DrawCam(screen, cam)
	p.bulletParticles.DrawCam(screen, cam)

}
func (p *Player) shoot() {
	if p.Ammo <= 0 {
		p.isShooting = false
		return
	}
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) || inpututil.IsKeyJustPressed(ebiten.KeyE) {
		p.isShooting = true
		p.Anim.ChangeAnim("shoot")
		bSpeed := 0
		//	p.shootParticles.Area.Y = int(p.Pos.Pos[1]) + 2
		if p.Anim.Flip {
			bSpeed = -5
			//		p.shootParticles.Area.X = p.Pos.Pos[0] - 16
		} else {
			//		p.shootParticles.Area.X = p.Pos.Pos[0] + 16
			bSpeed = 5
		}
		//	p.shootParticles.Spawn(10)
		p.Bullets = append(p.Bullets, NewBullet(p.Pos.Pos[0], p.Pos.Pos[1], bSpeed))
		p.Ammo--
	}
	if p.shootingTimer.Ticked() && p.isShooting {
		p.isShooting = false
	}
	if p.isShooting {
		p.shootingTimer.UpdateTimerTPS()
	}

}
func (p *Player) UpdateBullets(tiles map[[2]int]components.Rect, enemies []*Enemy) {
	for i, b := range p.Bullets {
		p.Bullets[i].Update(tiles, enemies)
		if b.Dead {
			p.bulletParticles.Area.X = b.Collider().X + b.Collider().Width
			p.bulletParticles.Area.Y = b.Collider().Y + b.Collider().Height/2
			p.bulletParticles.Spawn(10)
		}
	}
}
