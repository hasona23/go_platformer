package tilemap

import (
	"errors"
	"fmt"
	"go_platformer/components"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

type Level struct {
	Tilemap *TilemapJSON
	Objects *ObjectsJSON
}

func NewLevel(file []byte) *Level {
	var level Level
	var err error
	level.Tilemap, err = NewTilemapJSON(file)
	if err != nil {
		log.Fatal("Failed to Load Level Tilemap:", err)
	}

	level.Objects, err = NewObjectLayer(file)
	if err != nil {
		log.Fatal("Failed to Load Level Tilemap:", err)
	}
	return &level
}
func (l Level) GetCollisionTilesMap() map[[2]int]components.Rect {
	collisionMap := make(map[[2]int]components.Rect)
	for _, layer := range l.Tilemap.CollisionLayers {
		for _, tile := range layer.Tiles {
			collisionMap[[2]int{tile.X, tile.Y}] = components.NewRect(tile.X*l.Tilemap.TileSize, tile.Y*l.Tilemap.TileSize, l.Tilemap.TileSize, l.Tilemap.TileSize)
		}
	}
	return collisionMap
}
func (l Level) GetCollisionTiles() []components.Rect {
	r := []components.Rect{}
	for _, v := range l.GetCollisionTilesMap() {
		r = append(r, v)
	}
	return r
}
func (l Level) GetObjectsMap() map[string][]Object {
	return l.Objects.GetLayers()
}
func (l Level) GetObjectsByName(name string) ([]Object, error) {
	if _, ok := l.GetObjectsMap()[name]; ok {
		return l.GetObjectsMap()[name], nil
	}
	return nil, errors.New(fmt.Sprintln("couldnt find the objects of category name :", name))
}
func (l Level) GetSizeInPixels() [2]int {
	return [2]int{l.Tilemap.TileLayers[0].WidthInPixels, l.Tilemap.TileLayers[0].HeightInPixels}
}
func (l Level) Draw(screen *ebiten.Image, spriteSheet *ebiten.Image, showCollisions bool) {
	l.Tilemap.Draw(screen, spriteSheet, showCollisions)
}
func (l Level) DrawCamera(screen *ebiten.Image, spriteSheet *ebiten.Image, cam components.Camera, showCollisions bool) {
	l.Tilemap.DrawCamera(screen, spriteSheet, cam, showCollisions)
}
