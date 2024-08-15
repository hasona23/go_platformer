package entities

import (
	"go_platformer/components"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

type PhysicsEntity struct {
	Pos        components.Transform
	Anim       components.AnimSprite
	Vel        components.Vel
	Type       string
	Collisions map[string]bool
}

func NewPhyscisEntity(x, y, scale int, speed float32, etype string) *PhysicsEntity {
	e := &PhysicsEntity{}
	e.Type = etype
	e.Pos = components.Transform{Pos: [2]int{x, y}, Scale: float64(scale)}
	e.Vel.Speed = speed
	e.Collisions = map[string]bool{"right": false, "left": false, "up": false, "down": false}

	return e
}
func (e PhysicsEntity) Draw(screen *ebiten.Image, camera components.Camera) {
	op := &ebiten.DrawImageOptions{}
	e.Anim.Origin(op)
	if e.Anim.Flip {
		op.GeoM.Scale(-e.Pos.Scale, e.Pos.Scale)

	} else {
		op.GeoM.Scale(e.Pos.Scale, e.Pos.Scale)

	}

	op.GeoM.Translate(float64(e.Pos.Pos[0]+camera.X), float64(e.Pos.Pos[1]+camera.Y))

	screen.DrawImage(e.Anim.Img, op)
	//vector.DrawFilledRect(screen, float32(e.Collider().X+camera.X), float32(e.Collider().Y+camera.Y), float32(e.Collider().Width), float32(e.Collider().Height), color.Black, false)
}
func (e PhysicsEntity) GetAroundTiles(tiles map[[2]int]components.Rect) []components.Rect {
	offset := [][2]int{{1, 0}, {-1, 0}, {0, -1}, {0, 1}, {1, -1}, {-1, 1}, {-1, -1}, {1, 1}, {0, 0}}

	var tileRects []components.Rect
	for _, i := range offset {
		if tileRect, ok := tiles[[2]int{int(math.Ceil(float64(e.Collider().X)/float64(e.Anim.GetWidth()))) + i[0], int(math.Ceil(float64(e.Collider().Y)/float64(e.Anim.GetHeight()))) + i[1]}]; ok {
			tileRects = append(tileRects, tileRect)
		}

	}
	return tileRects
}
func (e PhysicsEntity) GetAroundTilesMap(tiles map[[2]int]components.Rect) map[[2]int]components.Rect {
	offset := [][2]int{{1, 0}, {-1, 0}, {0, -1}, {0, 1}, {1, -1}, {-1, 1}, {-1, -1}, {1, 1}, {0, 0}}

	var tileRects map[[2]int]components.Rect = make(map[[2]int]components.Rect)
	for _, i := range offset {
		if rect, ok := tiles[[2]int{int(math.Ceil(float64(e.Collider().X)/float64(e.Anim.GetWidth()))) + i[0], int(math.Ceil(float64(e.Collider().Y)/float64(e.Anim.GetHeight()))) + i[1]}]; ok {
			tileRects[[2]int{int(math.Ceil(float64(e.Collider().X)/float64(e.Anim.GetWidth()))) + i[0], int(math.Ceil(float64(e.Collider().Y)/float64(e.Anim.GetHeight()))) + i[1]}] = rect
		}

	}
	return tileRects
}
func (e *PhysicsEntity) Collider() components.Rect {
	bounds := e.Anim.Img.Bounds()
	return components.NewRect(
		int(e.Pos.Pos[0])-(bounds.Dx())/2,
		int(e.Pos.Pos[1])-(bounds.Dy())/2,
		int(bounds.Dx()),
		int(bounds.Dy()))
}

func (e *PhysicsEntity) Move(tiles map[[2]int]components.Rect) {
	//e.Vel.NormalizeDir()
	change := [2]int{int(math.Round(float64(e.Vel.Speed * e.Vel.Dir[0]))), int(math.Round(float64((e.Vel.Speed * e.Vel.Dir[1]))))}
	e.Collisions["up"] = false
	e.Collisions["down"] = false
	e.Collisions["left"] = false
	e.Collisions["right"] = false

	eRect := e.Collider()
	eRect.X += change[0]
	for _, rect := range e.GetAroundTiles(tiles) {
		if eRect.Collide(rect) {
			if change[0] > 0 {
				e.Collisions["right"] = true
			}
			if change[0] < 0 {
				e.Collisions["left"] = true
			}
		}
	}
	if !e.Collisions["right"] && change[0] > 0 {
		e.Pos.Pos[0] += change[0]
		e.Anim.Flip = false
	}
	if !e.Collisions["left"] && change[0] < 0 {
		e.Pos.Pos[0] += change[0]
		e.Anim.Flip = true
	}

	eRect = e.Collider()
	eRect.Y += change[1]
	for _, rect := range e.GetAroundTiles(tiles) {
		if eRect.Collide(rect) {
			if change[1] < 0 {
				e.Collisions["up"] = true
			}
			if change[1] > 0 {
				e.Collisions["down"] = true
			}
		}
	}

	if !e.Collisions["up"] && change[1] < 0 {
		e.Pos.Pos[1] += change[1]
	}
	if !e.Collisions["down"] && change[1] > 0 {

		e.Pos.Pos[1] += change[1]
	}
	if e.Collisions["up"] {
		e.Vel.Dir[1] *= -1
	}
	if e.Collisions["left"] {
		e.Anim.Flip = false
	}
	if e.Collisions["right"] {
		e.Anim.Flip = true
	}

	e.Anim.Update()

}
