package components

import "math"

type Vel struct {
	Dir   [2]float32
	Speed float32
}

func (v Vel) Length() float32 {
	return float32(math.Sqrt(float64(((v.Dir[0] * v.Dir[0]) + (v.Dir[1] * v.Dir[1])))))
}

func (v *Vel) NormalizeDir() {
	if v.Dir[0] != 0 && v.Dir[1] != 0 {
		factor := 1 / v.Length()
		v.Dir[0] *= factor
		v.Dir[1] *= factor
	}
}
func (v Vel) Init(speed float32) Vel {
	return Vel{[2]float32{0, 0}, speed}
}
