package particles

import (
	"go_platformer/spark"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

type ParticleOptsFunc func(*Particle)

type Particle struct {
	Img    *ebiten.Image
	X, Y   float32
	Scale  float32
	Raduis float32 //this is for circular motion
	Angle  float32
	Dir    spark.Vec2
	Speed  float32
	Color  color.Color
}

func DefaultParticle() Particle {
	return Particle{nil, 0, 0, 1, 0, 0, spark.Vec2{0, 0}, 0, color.RGBA{173, 216, 230, 255}}
}
func WithScale(scale float32) ParticleOptsFunc {
	return func(p *Particle) {
		p.Scale = scale
	}
}

// can be used  if no imae provided to render particles as squares of this color
func WithColor(color color.Color) ParticleOptsFunc {
	return func(p *Particle) {
		p.Color = color
	}
}
func WithImage(img *ebiten.Image) ParticleOptsFunc {
	return (func(p *Particle) {
		p.Img = img
	})
}
func WithRotation(raduis, speed float32) ParticleOptsFunc {
	return func(p *Particle) {
		p.Raduis = raduis
		p.Speed = speed
	}
}
func WithSpeed(speed float32) ParticleOptsFunc {
	return func(p *Particle) {
		p.Speed = speed
	}
}
func WithVelocity(dir spark.Vec2, speed float32) ParticleOptsFunc {
	return func(p *Particle) {
		p.Dir = dir
		p.Speed = speed
		p.Dir.NormalizeDir()
	}
}

// angle in degree
func WithAngle(angle float32) ParticleOptsFunc {
	return func(p *Particle) {
		p.Angle = angle
	}
}
func WithPos(x, y float32) ParticleOptsFunc {
	return func(p *Particle) {
		p.X = x
		p.Y = y

	}
}
func NewParticle(opts ...ParticleOptsFunc) *Particle {
	p := DefaultParticle()
	for _, fn := range opts {
		fn(&p)
	}
	return &p
}
