package components

import (
	"math"
)

type Camera struct {
	X, Y int
}

func NewCamera(x, y int) *Camera {
	return &Camera{
		X: x,
		Y: y,
	}
}

// ShakeAmount is how many times you want the shaking to occur
/*func (c *Camera) Shake(strengh float32, ShakeAmount int) {
	//factors is between  -1 , 1 to maniplutate direction of x and y displacement
	factor := 0
	for range ShakeAmount {
		if rand.Int()%2 == 0 {
			factor = 1
		} else {
			factor = -1
		}
		x := rand.Float32() * strengh * float32(factor)
		y := rand.Float32() * strengh * float32(factor)
		c.X += int(x)
		c.Y += int(y)
	}
}*/

// For Smooth Movement
func (c *Camera) FollowTarget(targetX, targetY, screenWidth, screenHeight, speed int) {
	c.X += int((-targetX + screenWidth/2.0 - c.X)) / speed
	c.Y += int((-targetY + screenHeight/2.0 - c.Y)) / speed
}

// For Sudden and fast without delay movemnt
func (c *Camera) GoTo(targetX, targetY, screenWidth, screenHeight int) {
	c.X = int((-targetX + screenWidth/2.0))
	c.Y = int((-targetY + screenHeight/2.0))
}

func (c *Camera) Constrain(tilemapWidthPixels, tilemapHeightPixels, screenWidth, screenHeight int) {
	c.X = int(math.Min(float64(c.X), 0.0))
	c.Y = int(math.Min(float64(c.Y), 0.0))

	c.X = int(math.Max(float64(c.X), float64(screenWidth-tilemapWidthPixels)))
	c.Y = int(math.Max(float64(c.Y), float64(screenHeight-tilemapHeightPixels)))
}
