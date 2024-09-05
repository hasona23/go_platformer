package particles

import (
	"go_platformer/components"
	"math"
	"math/rand"
	"slices"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type MotionType int

const (
	SingleDirection MotionType = iota
	Circular
	RandomDirections
	Inward
	Outward
)

type ParticleSystem struct {
	Particles          []Particle
	Motion             MotionType
	Area               components.Rect
	ModelParticle      Particle //this is the particle going to be spawned
	IsLooped           bool
	SpawnTime          components.Timer
	Decelration        float32
	ParticleSpawnCount uint
}

func (ps ParticleSystem) Raduis() float32 {
	return (float32(ps.Area.Width) + float32(ps.Area.Height)) / 2

}
func NewParticleSystem(x, y, width, height int, spawnRate float32, motion MotionType, particle Particle) *ParticleSystem {
	return &ParticleSystem{
		Area:               components.NewRect(x, y, width, height),
		Motion:             motion,
		ModelParticle:      particle,
		Decelration:        0.1,
		ParticleSpawnCount: 10,
	}
}
func (ps *ParticleSystem) Spawn(amount uint) {
	cX, cY := ps.Area.Centre()
	for range amount {
		x := float32(ps.Area.X) + (rand.Float32() * float32(ps.Area.Width))
		y := float32(ps.Area.Y) + (rand.Float32() * float32(ps.Area.Height))
		switch ps.Motion {
		case SingleDirection:
			ps.Particles = append(ps.Particles, *NewParticle(float64(x), float64(y), WithImage(ps.ModelParticle.Img), WithScale(ps.ModelParticle.Scale),
				WithVelocity(ps.ModelParticle.Dir, ps.ModelParticle.Speed)))
		case Circular:
			angle := rand.Float32() * 2 * math.Pi
			raduis := ps.Raduis() - rand.Float32()*ps.Raduis()/1.25
			rX := float64(cX) + float64(raduis)*math.Cos(float64(angle))
			rY := float64(cY) + float64(raduis)*math.Sin(float64(angle))

			ps.Particles = append(ps.Particles, *NewParticle(float64(rX), float64(rY),
				WithImage(ps.ModelParticle.Img), WithScale(ps.ModelParticle.Scale),
				WithRotation(float32(raduis), ps.ModelParticle.Speed), WithAngle(angle)))

		case RandomDirections:
			n1 := rand.Float32()
			n2 := rand.Float32()
			if n1 < 0.5 {
				n1 = -1
			}
			if n2 < .5 {
				n2 = -1
			}
			ps.Particles = append(ps.Particles, *NewParticle(float64(x), float64(y), WithImage(ps.ModelParticle.Img), WithScale(ps.ModelParticle.Scale),
				WithVelocity(Vec2f{rand.Float32() * n1, rand.Float32() * n2}, ps.ModelParticle.Speed)))
		case Inward:

			ps.Particles = append(ps.Particles, *NewParticle(float64(x), float64(y), WithImage(ps.ModelParticle.Img), WithScale(ps.ModelParticle.Scale),
				WithVelocity(Vec2f{float32(cX) - x, float32(cY) - y}, ps.ModelParticle.Speed)))
		case Outward:

			ps.Particles = append(ps.Particles, *NewParticle(float64(x), float64(y), WithImage(ps.ModelParticle.Img), WithScale(ps.ModelParticle.Scale),
				WithVelocity(Vec2f{x - float32(cX), y - float32(cY)}, ps.ModelParticle.Speed)))

		}
	}
}
func (ps *ParticleSystem) Update() {
	for i := range ps.Particles {
		ps.Particles[i].Speed -= ps.Decelration
	}
	ps.moveParticles()
	ps.Particles = slices.DeleteFunc(ps.Particles, func(p Particle) bool {
		return p.Speed <= 0
	})
	ps.SpawnTime.UpdateTimer()
	if ps.IsLooped && ps.SpawnTime.Ticked() {
		ps.Spawn(ps.ParticleSpawnCount)
	}
}
func (ps *ParticleSystem) moveParticles() {

	for i := range ps.Particles {
		switch ps.Motion {
		case Circular:
			ps.Particles[i].Angle += ps.Particles[i].Speed
			cX, cY := ps.Area.Centre()
			// NOTE: might use in future for sprial shapes or merging motions
			//ps.Particles[i].Raduis += ps.ModelParticle.Speed * 2
			//ps.Particles[i].Raduis = float32(math.Min(float64(ps.Particles[i].Raduis), float64(ps.Raduis())))
			ps.Particles[i].X = float64(cX) + float64(ps.Particles[i].Raduis)*math.Cos(float64(ps.Particles[i].Angle))
			ps.Particles[i].Y = float64(cY) + float64(ps.Particles[i].Raduis)*math.Sin(float64(ps.Particles[i].Angle))

		default:
			ps.Particles[i].X += float64(ps.Particles[i].Speed * ps.Particles[i].Dir[0])
			ps.Particles[i].Y += float64(ps.Particles[i].Speed * ps.Particles[i].Dir[1])

		}
	}
}
func (ps ParticleSystem) Draw(screen *ebiten.Image) {
	ps.drawWithOffset(screen, components.Point{X: 0, Y: 0})

}
func (ps ParticleSystem) DrawCam(screen *ebiten.Image, cam components.Camera) {
	ps.drawWithOffset(screen, components.Point(cam))
}
func (ps ParticleSystem) drawWithOffset(screen *ebiten.Image, offset components.Point) {
	if ps.ModelParticle.Img != nil {
		op := &ebiten.DrawImageOptions{}
		for _, p := range ps.Particles {
			p.Rotate(op, float64(math.Pi/2+p.Angle))
			op.GeoM.Translate(p.X/float64(p.Scale)+float64(offset.X), p.Y/float64(p.Scale)+float64(offset.Y))
			op.GeoM.Scale(float64(p.Scale), float64(p.Scale))
			screen.DrawImage(p.Img, op)
			op.GeoM.Reset()
		}
	} else {
		for _, p := range ps.Particles {
			vector.DrawFilledRect(screen, float32(p.X+float64(offset.X)), (float32(p.Y) + float32(offset.Y)), p.Scale, p.Scale, ps.ModelParticle.Color, false)
		}
	}
}
