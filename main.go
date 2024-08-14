package main

import (
	"go_platformer/assets"
	"go_platformer/components"
	"go_platformer/tilemap"
	"image"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	DisplayWidth  = 320
	DisplayHeight = 240
)

type Game struct {
	cam     components.Camera
	level1  *tilemap.Level
	enemies []tilemap.Object
}

func (g *Game) Init() {
	var err error
	g.level1 = tilemap.NewLevel(assets.Level1Map)
	g.cam = *components.NewCamera(0, 0)
	X = 0
	Y = 0
	g.enemies, err = g.level1.GetObjectsByName("Enemies")
	if err != nil {
		log.Fatal("Error Getting enemy objects :", err)
	}
}

var X int
var Y int

func (g *Game) Update() error {

	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		Y -= 10
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
		Y += 10
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		X -= 10
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		X += 20
	}
	g.cam.GoTo(X, Y, DisplayWidth, DisplayHeight)
	g.cam.Constrain(g.level1.GetSizeInPixels()[0], g.level1.GetSizeInPixels()[1], DisplayWidth, DisplayHeight)
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.level1.DrawCamera(screen, assets.SpriteSheet, g.cam, false)

	for _, object := range g.enemies {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(object.X), float64(object.Y))
		op.GeoM.Translate(float64(g.cam.X), float64(g.cam.Y))
		img := assets.SpriteSheet.SubImage(image.Rect(1*16, 5*16, 2*16, 6*16)).(*ebiten.Image)
		screen.DrawImage(img, op)
		op.GeoM.Reset()
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return DisplayWidth, DisplayHeight
}

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("GoPlatformer")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	assets.InitAssets()
	g := &Game{}
	g.Init()
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
