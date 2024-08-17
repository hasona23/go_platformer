package entities

import (
	"go_platformer/assets"
	"go_platformer/components"
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Player struct {
	*PhysicsEntity
	Bullets   []*Bullet
	isJumping bool
}

func NewPlayer() *Player {
	player := &Player{}
	player.PhysicsEntity = NewPhyscisEntity(40, 330, 1, 1, "player")
	player.Anim = components.NewAnimeSprite(assets.SpriteSheet, 16, 16)
	player.Anim.Add(components.NewAnimationFrame(0, 1, 4, 5, 2, "idle"))
	player.Anim.Add(components.NewAnimationFrame(0, 2, 4, 5, .2, "run"))
	player.Anim.Add(components.NewAnimationFrame(2, 3, 4, 5, 0.7, "shoot"))
	player.isJumping = false
	return player
}
func Lerp(a, b, t float32) float32 {
	return a + (b-a)*t
}
func (p *Player) Update(tiles map[[2]int]components.Rect) {
	//reached peak of jump
	if p.Collisions["down"] {
		p.Vel.Dir[1] = 0
		p.isJumping = false
		//fmt.Println(p.isJumping)
	}
	p.Vel.Dir[0] = 0
	if p.Vel.Dir[1] > 0 && p.isJumping {
		p.Vel.Dir[1] = min(3, p.Vel.Dir[1]+.8)
	} else {
		p.Vel.Dir[1] = min(2, p.Vel.Dir[1]+0.1)
	}
	if (inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) || inpututil.IsKeyJustPressed(ebiten.KeyW) || inpututil.IsKeyJustPressed(ebiten.KeySpace)) && !p.isJumping {
		p.Vel.Dir[1] = -3
		p.isJumping = true
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowDown) || ebiten.IsKeyPressed(ebiten.KeyS) {
		p.Vel.Dir[1] = 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) || ebiten.IsKeyPressed(ebiten.KeyA) {
		p.Vel.Dir[0] = -1
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) || ebiten.IsKeyPressed(ebiten.KeyD) {
		p.Vel.Dir[0] = 1
	}

	if p.Vel.Dir[0] != 0 {
		p.Anim.ChangeAnim("run")
	} else {
		p.Anim.ChangeAnim("idle")
	}
	p.PhysicsEntity.Move(tiles)
	p.shoot()

}
func (p Player) Draw(screen *ebiten.Image, cam components.Camera) {
	for _, b := range p.Bullets {
		b.Draw(screen, cam)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) || ebiten.IsKeyPressed(ebiten.KeyE) {
		op := &ebiten.DrawImageOptions{}
		if !p.Anim.Flip {
			op.GeoM.Translate(float64(p.Pos.Pos[0]+cam.X+8), float64(float64(p.Pos.Pos[1]+cam.Y)-5.5))
		} else {

			p.Anim.Origin(op)
			op.GeoM.Scale(-1, 1)
			op.GeoM.Translate(float64(p.Pos.Pos[0]+cam.X-16), float64(float64(p.Pos.Pos[1]+cam.Y))+2)

		}
		screen.DrawImage(assets.SpriteSheet.SubImage(image.Rect(3*16, 4*16, 4*16, 5*16)).(*ebiten.Image), op)
	}
	p.PhysicsEntity.Draw(screen, cam)

}
func (p *Player) shoot() {
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) || inpututil.IsKeyJustPressed(ebiten.KeyE) {
		bSpeed := 0
		if p.Anim.Flip {
			bSpeed = -5
		} else {
			bSpeed = 5
		}
		p.Bullets = append(p.Bullets, NewBullet(p.Pos.Pos[0], p.Pos.Pos[1], bSpeed))
		p.Anim.ChangeAnim("shoot")
	}
}
func (p *Player) UpdateBullets(tiles map[[2]int]components.Rect, enemies []*Enemy) {
	for i := range p.Bullets {
		p.Bullets[i].Update(tiles, enemies)
	}
}
