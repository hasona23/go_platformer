package tilemap

import (
	"errors"
	"fmt"

	"log"

	"go_platformer/spark"

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
func (l Level) GetCollisionTilesMap() map[[2]int]spark.Rect {
	collisionMap := make(map[[2]int]spark.Rect)
	for _, layer := range l.Tilemap.CollisionLayers {
		for _, tile := range layer.Tiles {
			collisionMap[[2]int{tile.X, tile.Y}] = spark.NewRect(tile.X*l.Tilemap.TileSize, tile.Y*l.Tilemap.TileSize, l.Tilemap.TileSize, l.Tilemap.TileSize)
		}
	}
	return collisionMap
}
func (l Level) GetCollisionTiles() []spark.Rect {
	r := []spark.Rect{}
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
func (l Level) DrawCamera(screen *ebiten.Image, spriteSheet *ebiten.Image, cam spark.Cam, showCollisions bool) {
	l.Tilemap.DrawCamera(screen, spriteSheet, cam, showCollisions)
}
