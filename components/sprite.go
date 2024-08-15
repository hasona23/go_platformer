package components

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Sprite struct {
	Img  *ebiten.Image
	Flip bool
}

func (sprite Sprite) GetWidth() int {
	return sprite.Img.Bounds().Dx()
}

// use for rotation ,scaling , reversing face/inverse scale
func (sprite *Sprite) Origin(op *ebiten.DrawImageOptions) {
	s := sprite.Img.Bounds()
	op.GeoM.Translate(-float64(s.Dx())/2, -float64(s.Dy())/2)

}
func (sprite *Sprite) Rotate(op *ebiten.DrawImageOptions, angle float64) {
	s := sprite.Img.Bounds()
	op.GeoM.Translate(-float64(s.Dx())/2, -float64(s.Dy())/2)
	op.GeoM.Rotate(angle)
	op.GeoM.Translate(float64(s.Dx())/2, float64(s.Dy())/2)
}
