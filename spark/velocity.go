package spark

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

// returns the x velocity
func (v Vel) VelX() float32 {
	return v.Speed * v.Dir[0]
}

func (v Vel) VelY() float32 {
	return v.Speed * v.Dir[1]
}
