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
	//	e.Anim.Origin(op)
	if e.Anim.Flip {
		e.Anim.FlipHorizontal(op)
	}
	op.GeoM.Translate(float64(e.Pos.Pos[0]+camera.X), float64(e.Pos.Pos[1]+camera.Y))

	screen.DrawImage(e.Anim.Img, op)
	// vector.StrokeRect(screen, float32(e.Collider().X+camera.X), float32(e.Collider().Y+camera.Y),
	// float32(e.Collider().Width), float32(e.Collider().Height), 2, color.Black, false)
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

// Gets from the map of collision tileset from tilemap
//
// this gets the collision tiles in all directions of player (square)
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

// Collider return the entity rect
//
// Collider is adjusted so that pos is moved half size of sprite up and to left
func (e *PhysicsEntity) Collider() components.Rect {
	bounds := e.Anim.Img.Bounds()
	return components.NewRect(
		int(e.Pos.Pos[0]),
		int(e.Pos.Pos[1]),
		int(bounds.Dx()),
		int(bounds.Dy()))
}

// Move Moves the entity with collision WithTiles
const (
	Up    = "up"
	Down  = "down"
	Left  = "left"
	Right = "right"
)

func (e *PhysicsEntity) Move(tiles map[[2]int]components.Rect) {
	change := e.calculateChange()
	e.resetCollisions()

	e.handleHorizontalMovement(tiles, change[0])
	e.handleVerticalMovement(tiles, change[1])

	e.updateVelocity()
	e.updateAnimation()
}

func (e *PhysicsEntity) calculateChange() [2]int {
	return [2]int{
		int(math.Round(float64(e.Vel.Speed * e.Vel.Dir[0]))),
		int(math.Round(float64(e.Vel.Speed * e.Vel.Dir[1]))),
	}
}

func (e *PhysicsEntity) resetCollisions() {
	e.Collisions = map[string]bool{
		Up: false, Down: false, Left: false, Right: false,
	}
}

func (e *PhysicsEntity) handleHorizontalMovement(tiles map[[2]int]components.Rect, changeX int) {
	eRect := e.Collider()
	eRect.X += changeX

	for _, rect := range e.GetAroundTiles(tiles) {
		if eRect.Collide(rect) {
			if changeX > 0 {
				e.Collisions[Right] = true
			} else if changeX < 0 {
				e.Collisions[Left] = true
			}
			return
		}
	}

	if !e.Collisions[Left] && changeX < 0 || !e.Collisions[Right] && changeX > 0 {
		e.Pos.Pos[0] += changeX
		e.Anim.Flip = changeX < 0
	}
}

func (e *PhysicsEntity) handleVerticalMovement(tiles map[[2]int]components.Rect, changeY int) {
	eRect := e.Collider()
	eRect.Y += changeY

	for _, rect := range e.GetAroundTiles(tiles) {
		if eRect.Collide(rect) {
			if changeY < 0 {
				e.Collisions[Up] = true
			} else if changeY > 0 {
				e.Collisions[Down] = true
			}
			return
		}
	}

	if !e.Collisions[Up] && changeY < 0 || !e.Collisions[Down] && changeY > 0 {
		e.Pos.Pos[1] += changeY
	}
}

func (e *PhysicsEntity) updateVelocity() {
	if e.Collisions[Up] && e.Vel.Dir[1] < 0 {
		e.Vel.Dir[1] = -0.25
	}

}

func (e *PhysicsEntity) updateAnimation() {
	if e.Collisions[Left] {
		e.Anim.Flip = false
	} else if e.Collisions[Right] {
		e.Anim.Flip = true
	}
	e.Anim.Update()
}
