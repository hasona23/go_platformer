package spark

import "math"

type Vec2 [2]float32

func (v Vec2) Length() float32 {
	return float32(math.Sqrt(float64(((v[0] * v[0]) + (v[1] * v[1])))))
}

func (v *Vec2) NormalizeDir() {
	if v[0] != 0 && v[1] != 0 {
		factor := 1 / v.Length()
		v[0] *= factor
		v[1] *= factor
	}
}
