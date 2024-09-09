package spark

import "github.com/hajimehoshi/ebiten/v2"

type SpriteEffect uint // Flip horizontal and vertical

const (
	None SpriteEffect = iota
	FlipHorizontal
	FlipVertical
)

// an interface for sprite and animated sprite
type Spriter interface {
	GetHeight() int
	GetWidth() int
	SetSpriteOP(op *ebiten.DrawImageOptions, rotation float32)
}
