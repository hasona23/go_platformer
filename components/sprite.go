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
func (sprite Sprite) GetHeight() int {
	return sprite.Img.Bounds().Dy()
}

func (sprite Sprite) FlipHorizontal(op *ebiten.DrawImageOptions) {
	s := sprite.Img.Bounds()
	op.GeoM.Scale(-1, 1)
	//op.GeoM.Translate(float64(s.Bounds().Dx())/2, 0)
	op.GeoM.Translate(float64(s.Dx()), 0)

}
func (sprite Sprite) FlipVertical(op *ebiten.DrawImageOptions) {
	s := sprite.Img.Bounds()
	op.GeoM.Scale(1, -1)
	//op.GeoM.Translate(float64(s.Bounds().Dx())/2, 0)
	op.GeoM.Translate(0, float64(s.Dy()))

}
func (sprite *Sprite) Rotate(op *ebiten.DrawImageOptions, angle float64) {
	s := sprite.Img.Bounds()
	op.GeoM.Translate(-float64(s.Dx())/2, -float64(s.Dy())/2)
	op.GeoM.Rotate(angle)
	op.GeoM.Translate(float64(s.Dx())/2, float64(s.Dy())/2)
}
