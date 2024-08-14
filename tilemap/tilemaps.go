package tilemap

import (
	"encoding/json"
	"go_platformer/components"
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Tile struct {
	X, Y, Id int
}
type TilemapLayerJSON struct {
	Data           []int  `json:"data"`
	Height         int    `json:"height"`
	Width          int    `json:"width"`
	Name           string `json:"name"`
	Class          string `json:"class"`
	Tiles          []Tile
	TileSize       int
	WidthInPixels  int
	HeightInPixels int
}

type TilemapJSON struct {
	TileLayers      []TilemapLayerJSON `json:"layers"`
	TileSize        int                `json:"tilewidth"`
	CollisionLayers []TilemapLayerJSON
}

func NewTilemapJSON(file []byte) (*TilemapJSON, error) {
	contents := file
	var tilemap TilemapJSON
	err := json.Unmarshal(contents, &tilemap)
	for i, layer := range tilemap.TileLayers {
		tilemap.TileLayers[i].TileSize = tilemap.TileSize
		tilemap.TileLayers[i].WidthInPixels = tilemap.TileSize * tilemap.TileLayers[i].Width
		tilemap.TileLayers[i].HeightInPixels = tilemap.TileSize * tilemap.TileLayers[i].Height
		for j, Id := range tilemap.TileLayers[i].Data {
			if Id == 0 {
				continue
			}
			tile := Tile{}
			tile.Id = Id
			tile.X = (j % layer.Width)
			tile.Y = (j / layer.Width)
			tilemap.TileLayers[i].Tiles = append(tilemap.TileLayers[i].Tiles, tile)

		}

		if layer.Class == "collision" {
			tilemap.CollisionLayers = append(tilemap.CollisionLayers, tilemap.TileLayers[i])

		}
	}

	if err != nil {
		return nil, err
	}
	return &tilemap, err
}

func (t *TilemapJSON) Draw(screen *ebiten.Image, spriteSheet *ebiten.Image, showCollision bool) {
	for _, layer := range t.TileLayers {
		for _, tile := range layer.Tiles {
			tilesInRow := spriteSheet.Bounds().Dx() / layer.TileSize
			srcX := (tile.Id - 1) % tilesInRow
			srcY := (tile.Id - 1) / tilesInRow
			srcX *= layer.TileSize
			srcY *= layer.TileSize

			op := &ebiten.DrawImageOptions{}

			op.GeoM.Translate(float64(tile.X*layer.TileSize), float64(tile.Y*layer.TileSize))

			screen.DrawImage(spriteSheet.SubImage(image.Rect(srcX, srcY, srcX+layer.TileSize, srcY+layer.TileSize)).(*ebiten.Image), op)
			op.GeoM.Reset()
			if showCollision && layer.Class == "collision" {
				vector.DrawFilledRect(screen, float32(tile.X*layer.TileSize), float32(tile.Y*layer.TileSize), float32(layer.TileSize),
					float32(layer.TileSize), color.RGBA{123, 123, 123, 123}, false)

			}
		}
	}

}

func (t *TilemapJSON) DrawCamera(screen *ebiten.Image, spriteSheet *ebiten.Image, camera components.Camera, showCollision bool) {
	for _, layer := range t.TileLayers {
		for _, tile := range layer.Tiles {
			if tile.Id == 0 {
				continue
			}
			tilesInRow := spriteSheet.Bounds().Dx() / layer.TileSize
			srcX := (tile.Id - 1) % tilesInRow
			srcY := (tile.Id - 1) / tilesInRow
			srcX *= layer.TileSize
			srcY *= layer.TileSize

			op := &ebiten.DrawImageOptions{}

			op.GeoM.Translate(float64(tile.X*layer.TileSize), float64(tile.Y*layer.TileSize))
			op.GeoM.Translate(float64(camera.X), float64(camera.Y))
			screen.DrawImage(spriteSheet.SubImage(image.Rect(srcX, srcY, srcX+layer.TileSize, srcY+layer.TileSize)).(*ebiten.Image), op)
			op.GeoM.Reset()
			if showCollision && layer.Class == "collision" {
				//fmt.Println(float32(tile.X*layer.TileSize+camera.X), float32(tile.Y*layer.TileSize+camera.Y))
				vector.DrawFilledRect(screen, float32(tile.X*layer.TileSize+camera.X), float32(tile.Y*layer.TileSize+camera.Y), float32(layer.TileSize),
					float32(layer.TileSize), color.RGBA{123, 123, 123, 123}, false)

			}
		}

	}

}
