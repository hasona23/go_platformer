package assets

import (
	"bytes"
	"embed"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

var (
	//go:embed sprites/*
	spritesFiles embed.FS
	//go:embed sound/*
	soundFiles embed.FS
	//go:embed maps/*
	mapFiles embed.FS
	//go:embed fonts/*
	fontFiles embed.FS
)
var (
	SpriteSheet *ebiten.Image
	Level1Map   []byte
	PixelFont   []byte
)

func InitAssets() {
	var err error
	spritesheetFile, err := spritesFiles.ReadFile("sprites/tilemap_packed.png")
	if err != nil {
		log.Fatal("Error Loading Assets :", err)
	}
	SpriteSheet, _, err = ebitenutil.NewImageFromReader(bytes.NewReader(spritesheetFile))
	if err != nil {
		log.Fatal("Error Loading Assets :", err)
	}
	Level1Map, err = mapFiles.ReadFile("maps/level1.json")
	if err != nil {
		log.Fatal("Error Loading Level 1")
	}
	PixelFont, err = fontFiles.ReadFile("fonts/pixelFont.ttf")
	if err != nil {
		log.Fatal("Error loading PixelFont :", err)
	}
}
