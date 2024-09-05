package particles

import (
	"go_platformer/components"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

type Vec2f [2]float32

func (v *Vec2f) NormalizeDir() {
	if v[0] != 0 && v[1] != 0 {
		factor := 1 / v.Length()
		v[0] *= factor
		v[1] *= factor
	}
}
func (v Vec2f) Length() float32 {
	return float32(math.Sqrt(float64(((v[0] * v[0]) + (v[1] * v[1])))))
}

type ParticleOptsFunc func(*ParticleOptions)
type ParticleOptions struct {
	components.Sprite
	X, Y   float64
	Scale  float32
	Raduis float32 //this is for circular motion
	Angle  float32
	Dir    Vec2f
	Speed  float32
	Color  color.Color
}

type Particle struct {
	ParticleOptions
}

func DefualtOpts(x, y float64) ParticleOptions {
	return ParticleOptions{components.Sprite{}, x, y, 1, 0, 0, Vec2f{0, 0}, 0, color.RGBA{173, 216, 230, 255}}
}
func WithScale(scale float32) ParticleOptsFunc {
	return func(po *ParticleOptions) {
		po.Scale = scale
	}
}

// can be used  if no imae provided to render particles as squares of this color
func WithColor(color color.Color) ParticleOptsFunc {
	return func(po *ParticleOptions) {
		po.Color = color
	}
}
func WithImage(img *ebiten.Image) ParticleOptsFunc {
	return (func(po *ParticleOptions) {
		po.Img = img
	})
}
func WithRotation(raduis, speed float32) ParticleOptsFunc {
	return func(po *ParticleOptions) {
		po.Raduis = raduis
		po.Speed = speed
	}
}
func WithSpeed(speed float32) ParticleOptsFunc {
	return func(po *ParticleOptions) {
		po.Speed = speed
	}
}
func WithVelocity(dir Vec2f, speed float32) ParticleOptsFunc {
	return func(po *ParticleOptions) {
		po.Dir = dir
		po.Speed = speed
		po.Dir.NormalizeDir()
	}
}

// angle in degree
func WithAngle(angle float32) ParticleOptsFunc {
	return func(po *ParticleOptions) {
		po.Angle = angle
	}
}

func NewParticle(x, y float64, opts ...ParticleOptsFunc) *Particle {
	po := DefualtOpts(x, y)
	for _, fn := range opts {
		fn(&po)
	}
	return &Particle{po}
}
