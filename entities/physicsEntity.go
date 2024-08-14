package entities

import (
	"go_platformer/components"
	"math"
)

type PhysicsEntity struct {
	components.Transform
	components.AnimSprite
	components.Vel
}
type PhysicsEntityManager struct {
	entities []*PhysicsEntity
	Tiles    map[[2]int]components.Rect
	offset   [][2]int
}

func NewPhyscisEntity(x, y, scale int) PhysicsEntity {
	e := PhysicsEntity{}
	e.Transform = components.Transform{Pos: [2]int{x, y}, Scale: float64(scale)}

	return e
}
func (e PhysicsEntity) GetAroundTiles(tiles map[[2]int]components.Rect) []components.Rect {
	offset := [][2]int{{1, 0}, {-1, 0}, {0, -1}, {0, 1}, {1, -1}, {-1, 1}, {-1, -1}, {1, 1}, {0, 0}}

	var tileRects []components.Rect
	for _, i := range offset {
		if tileRect, ok := tiles[[2]int{int(math.Ceil(float64(e.Collider().X/8))) + i[0], int(math.Ceil(float64(e.Collider().Y/8))) + i[1]}]; ok {
			tileRects = append(tileRects, tileRect)
		}

	}
	return tileRects
}
func (pe *PhysicsEntity) Collider() components.Rect {
	bounds := pe.AnimSprite.Img.Bounds()
	return components.NewRect(
		int(pe.Transform.Pos[0])-(bounds.Dx()),
		int(pe.Transform.Pos[1])-(bounds.Dy()),
		int(bounds.Dx()),
		int(bounds.Dy()))
}
func (e *PhysicsEntity) Move(tiles map[[2]int]components.Rect, dt float32) {
	change := [2]int{int(e.Vel.Speed * e.Vel.Dir[0] * dt), int(e.Vel.Speed * e.Vel.Dir[1] * dt)}
	// 					  right   left   up   down
	isColliding := [4]bool{false, false, false, false}

	eRect := e.Collider()
	eRect.X += change[0]
	for _, rect := range e.GetAroundTiles(tiles) {
		if eRect.Collide(rect) {
			if change[0] > 0 {
				isColliding[0] = true
			}
			if change[0] < 0 {
				isColliding[1] = true
			}
		}
	}
	if !isColliding[0] && change[0] > 0 {
		e.Transform.Pos[0] += change[0]
	}
	if !isColliding[1] && change[0] < 0 {
		e.Transform.Pos[0] += change[0]
	}

	eRect = e.Collider()
	eRect.Y += change[1]
	for _, rect := range e.GetAroundTiles(tiles) {
		if eRect.Collide(rect) {
			if change[1] > 0 {
				isColliding[3] = true
			}
			if change[1] < 0 {
				isColliding[2] = true
			}
		}
	}
	if !isColliding[3] && change[1] > 0 {
		e.Transform.Pos[1] += change[1]
	}
	if !isColliding[2] && change[1] < 0 {
		e.Transform.Pos[1] += change[1]
	}

}
