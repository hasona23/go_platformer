package main

import (
	"go_platformer/assets"
	"go_platformer/components"
	"go_platformer/entities"
	"go_platformer/tilemap"
	ui "go_platformer/ui"
	"image/color"
	"log"
	"slices"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	DisplayWidth  = 320
	DisplayHeight = 240
)

type Game struct {
	cam       components.Camera
	level1    *tilemap.Level
	enemies   []*entities.Enemy
	player    *entities.Player
	label     *ui.Label
	state     GameState
	uimanager *ui.UIManager
}

func (g *Game) Init() {
	var err error
	g.level1 = tilemap.NewLevel(assets.Level1Map)
	g.cam = *components.NewCamera(0, 0)
	enemyObjects, err := g.level1.GetObjectsByName("Enemies")
	for _, i := range enemyObjects {
		g.enemies = append(g.enemies, entities.NewEnemy(int(i.X), int(i.Y)))
	}
	if err != nil {
		log.Fatal("Error Getting enemy objects :", err)
	}
	g.state = Start
	g.player = entities.NewPlayer()

	g.label = ui.NewLabel("Hello", 5, 5, assets.PixelFont, 16, color.Black)
	button1 := ui.NewButton("start", 100, 100, 16, 2, assets.PixelFont, color.Black, color.White)
	button2 := ui.NewButton("save", 100, 150, 16, 2, assets.PixelFont, color.Black, color.White)
	button3 := ui.NewButton("quit", 100, 200, 16, 2, assets.PixelFont, color.Black, color.White)
	g.uimanager = ui.NewUIManager()
	g.uimanager.AddButton(button1)
	g.uimanager.AddButton(button2)
	g.uimanager.AddButton(button3)

}

func (g *Game) Update() error {

	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		g.state = Pause
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) && g.state == Pause {
		g.state = Main
	}
	if g.state == Start {
		g.uimanager.Update()
	}
	if g.state == Main {
		if g.player.Died {
			g.enemies = []*entities.Enemy{}
			g.Init()
		}

		tilesCollisionMap := g.level1.GetCollisionTilesMap()
		g.player.Update(tilesCollisionMap)
		g.player.UpdateBullets(tilesCollisionMap, g.enemies)
		g.cam.FollowTarget(g.player.Pos.Pos[0], g.player.Pos.Pos[1], DisplayWidth, DisplayHeight, 30)
		g.cam.Constrain(g.level1.GetSizeInPixels()[0], g.level1.GetSizeInPixels()[1], DisplayWidth, DisplayHeight)
		for _, i := range g.enemies {
			if !i.Dead {
				i.Update(tilesCollisionMap, g.player)
			}
		}

		g.player.Bullets = slices.DeleteFunc(g.player.Bullets, func(b *entities.Bullet) bool { return b.Dead })
		g.enemies = slices.DeleteFunc(g.enemies, func(e *entities.Enemy) bool { return e.Dead })
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	if g.state != Start {
		g.level1.DrawCamera(screen, assets.SpriteSheet, g.cam, false)
		for _, object := range g.enemies {
			object.Draw(screen, g.cam)
		}
		g.player.Draw(screen, g.cam)
		g.player.PhysicsEntity.Draw(screen, g.cam)
	}
	//g.label.Text = fmt.Sprintf("Ammo:%d", g.player.Ammo)
	//g.label.Draw(screen)
	if g.state == Start {
		g.uimanager.Draw(screen)
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
