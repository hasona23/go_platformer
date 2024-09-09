package particles

import (
	"fmt"
	"go_platformer/spark"
	"math"
	"math/rand"
	"slices"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type MotionType int

type PSOptsFunc func(*ParticleSystem)

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
	Area               spark.Rect
	ModelParticle      Particle //this is the particle going to be spawned
	IsLooped           bool
	SpawnTime          spark.Timer
	Decelration        float32 //Decreases speed. when particles speed is zero particle dies
	Shrink             float32 //Decreases scale . When scale is zero particle dies
	Gravity            float32 // affect the Y velocity
	ParticleSpawnCount uint
}

func (ps ParticleSystem) Raduis() float32 {
	return (float32(ps.Area.Width) + float32(ps.Area.Height)) / 2

}

// default particle system
func DefaultPS() ParticleSystem {
	return ParticleSystem{make([]Particle, 64), Outward, spark.NewRect(0, 0, 16, 16), DefaultParticle(), false, spark.NewTimer(0), 0.0, 0, 0, 0}
}

// x,y,width,heigt the area where particles spawn
func WithArea(area spark.Rect) PSOptsFunc {
	return func(ps *ParticleSystem) {
		ps.Area = area
	}
}
func WithShrinking(rate float32) PSOptsFunc {
	return func(ps *ParticleSystem) {
		ps.Shrink = rate
	}
}
func WithGravity(strength float32) PSOptsFunc {
	return func(ps *ParticleSystem) {
		ps.Gravity = strength
	}
}
func WithLooping() PSOptsFunc {
	return func(ps *ParticleSystem) {
		ps.IsLooped = true
	}
}

// particles spawned per time
func WithParticleSpawnCount(n uint) PSOptsFunc {
	return func(ps *ParticleSystem) {
		ps.ParticleSpawnCount = n
	}
}

// make the particle system player looped
//
// makes the particle system play every spawnrate passes in seconds
func WithSpawnRate(rate float32) PSOptsFunc {
	return func(ps *ParticleSystem) {
		ps.IsLooped = true
		ps.SpawnTime = spark.NewTimer(rate)
	}
}
func WithModelParticle(particle Particle) PSOptsFunc {
	return func(ps *ParticleSystem) {
		ps.ModelParticle = particle
	}
}

// default is 0.1
func WithDecelration(decelration float32) PSOptsFunc {
	return func(ps *ParticleSystem) {
		ps.Decelration = decelration
	}
}
func WithMotionType(motionType MotionType) PSOptsFunc {
	return func(ps *ParticleSystem) {
		ps.Motion = motionType
	}
}
func NewParticleSystem(opts ...PSOptsFunc) *ParticleSystem {
	ps := DefaultPS()
	for _, fn := range opts {
		fn(&ps)
	}
	return &ps
}
func (ps *ParticleSystem) Spawn(amount uint) {
	cX, cY := ps.Area.Centre()
	for range amount {
		x := float32(ps.Area.X) + (rand.Float32() * float32(ps.Area.Width))
		y := float32(ps.Area.Y) + (rand.Float32() * float32(ps.Area.Height))
		switch ps.Motion {
		case SingleDirection:
			ps.Particles = append(ps.Particles, *NewParticle(WithPos(x, y), WithImage(ps.ModelParticle.Img), WithScale(ps.ModelParticle.Scale),
				WithVelocity(ps.ModelParticle.Dir, ps.ModelParticle.Speed)))
		case Circular:
			angle := rand.Float32() * 2 * math.Pi
			raduis := ps.Raduis() - rand.Float32()*ps.Raduis()
			ps.Particles = append(ps.Particles, *NewParticle(WithPos(float32(cY), float32(cX)),
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
			ps.Particles = append(ps.Particles, *NewParticle(WithPos(x, y), WithImage(ps.ModelParticle.Img), WithScale(ps.ModelParticle.Scale),
				WithVelocity(spark.Vec2{rand.Float32() * n1, rand.Float32() * n2}, ps.ModelParticle.Speed)))
		case Inward:

			ps.Particles = append(ps.Particles, *NewParticle(WithPos(x, y), WithImage(ps.ModelParticle.Img), WithScale(ps.ModelParticle.Scale),
				WithVelocity(spark.Vec2{float32(cX) - x, float32(cY) - y}, ps.ModelParticle.Speed)))
		case Outward:

			ps.Particles = append(ps.Particles, *NewParticle(WithPos(x, y), WithImage(ps.ModelParticle.Img), WithScale(ps.ModelParticle.Scale),
				WithVelocity(spark.Vec2{x - float32(cX), y - float32(cY)}, ps.ModelParticle.Speed)))
		}
	}
}
func (ps *ParticleSystem) Update() {
	for i := range ps.Particles {
		ps.Particles[i].Speed -= ps.Decelration
		ps.Particles[i].Scale -= ps.Shrink
		ps.Particles[i].Dir[1] += ps.Gravity * 0.1 //factor to avoid dealing with small flaoting point number
	}
	ps.moveParticles()
	ps.Particles = slices.DeleteFunc(ps.Particles, func(p Particle) bool {
		return p.Speed <= 0 || p.Scale <= 0
	})
	if ps.IsLooped {
		ps.SpawnTime.UpdateTimerTPS()
		if ps.SpawnTime.Ticked() {
			ps.Spawn(ps.ParticleSpawnCount)
		}
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
			ps.Particles[i].X = float32(cX) + (ps.Particles[i].Raduis)*float32(math.Cos(float64(ps.Particles[i].Angle)))
			ps.Particles[i].Y = float32(cY) + (ps.Particles[i].Raduis)*float32(math.Sin(float64(ps.Particles[i].Angle)))

		default:
			ps.Particles[i].X += (ps.Particles[i].Speed * ps.Particles[i].Dir[0])
			ps.Particles[i].Y += (ps.Particles[i].Speed * ps.Particles[i].Dir[1])

		}
	}
}
func (ps ParticleSystem) Draw(screen *ebiten.Image) {
	ps.drawWithOffset(screen, spark.Vec2{0, 0})

}
func (ps ParticleSystem) DrawCam(screen *ebiten.Image, cam spark.Cam) {
	ps.drawWithOffset(screen, spark.Vec2{cam.X, cam.Y})
}
func (ps ParticleSystem) drawWithOffset(screen *ebiten.Image, offset spark.Vec2) {
	if ps.ModelParticle.Img != nil {
		op := &ebiten.DrawImageOptions{}
		for _, p := range ps.Particles {
			s := ps.ModelParticle.Img.Bounds()
			op.GeoM.Translate(-float64(s.Dx())/2, -float64(s.Dy())/2)
			op.GeoM.Rotate(float64(p.Angle))
			op.GeoM.Translate(float64(p.X/p.Scale)+float64(offset[0]), float64(p.Y/float32(p.Scale)+float32(offset[1])))
			op.GeoM.Scale(float64(p.Scale), float64(p.Scale))
			screen.DrawImage(p.Img, op)
			op.GeoM.Reset()
			fmt.Println("img")
		}
	} else {
		for _, p := range ps.Particles {
			vector.DrawFilledRect(screen, (p.X + float32(offset[0])), (p.Y + float32(offset[1])), p.Scale, p.Scale, ps.ModelParticle.Color, false)
		}
	}
	//vector.StrokeRect(screen, float32(ps.Area.X)+offset[0], float32(ps.Area.Y)+offset[1], float32(ps.Area.Width), float32(ps.Area.Height), 10, color.Black, false)
}
