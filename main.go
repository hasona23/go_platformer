package main

import (
	"fmt"
	"go_platformer/assets"
	"go_platformer/components"
	"go_platformer/entities"
	"go_platformer/tilemap"
	ui "go_platformer/ui"
	"image/color"
	"log"
	"os"
	"slices"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	DisplayWidth  = 320
	DisplayHeight = 240
)

// uiComponents
const (
	Start       = "start"
	Save        = "save"
	Exit        = "exit"
	AmmoCounter = "ammo"
)

type Game struct {
	cam      components.Camera
	level1   *tilemap.Level
	enemies  []*entities.Enemy
	player   *entities.Player
	state    GameState
	mainmenu *ui.UILayout
	gameUI   *ui.UILayout
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
	g.state = MainMenu
	g.player = entities.NewPlayer()

	ammoUI := ui.NewLabel("Ammo:", 5, 5, assets.PixelFont, 16, color.Black)
	g.gameUI = ui.NewUILayout("MainGameUI")
	g.gameUI.AddLabel(AmmoCounter, ammoUI)
	start := ui.NewButton("start", 100, 100, 16, 2, assets.PixelFont, color.White, color.RGBA{255, 222, 206, 255}, color.White)
	save := ui.NewButton("save", 100, 150, 16, 2, assets.PixelFont, color.White, color.RGBA{255, 222, 206, 255}, color.White)
	exit := ui.NewButton("quit", 100, 200, 16, 2, assets.PixelFont, color.White, color.RGBA{255, 222, 206, 255}, color.White)
	start.Style.BorderThickness = 3
	save.Style.BorderThickness = 3
	exit.Style.BorderThickness = 3
	g.mainmenu = ui.NewUILayout("MainMenu")
	g.mainmenu.AddButton(Start, start)
	g.mainmenu.AddButton(Save, save)
	g.mainmenu.AddButton(Exit, exit)
	g.mainmenu.ApplyHoverToAllButtons(hover)
	start.OnClick = func(b *ui.Button) {
		if g.state != Main {
			g.state = Main
		}
	}

	exit.OnClick = func(b *ui.Button) {
		os.Exit(0)
	}
}

func hover(b *ui.Button) {
	b.Text.Style.Color = color.Black
	b.Style.BorderColor = color.White
}

func (g *Game) Update() error {

	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		g.state = Pause
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) && g.state == Pause {
		g.state = Main
	}
	if g.state == MainMenu {
		g.mainmenu.Update()
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
	if g.state == Main {
		g.level1.DrawCamera(screen, assets.SpriteSheet, g.cam, false)
		for _, object := range g.enemies {
			object.Draw(screen, g.cam)
		}
		g.player.Draw(screen, g.cam)
		g.player.PhysicsEntity.Draw(screen, g.cam)
		counter, _ := g.gameUI.GetLabel(AmmoCounter)
		counter.Text = fmt.Sprintf("Ammo:%d", g.player.Ammo)
		g.gameUI.Draw(screen)
	}

	if g.state == MainMenu {
		screen.Fill(color.RGBA{255, 222, 206, 255})
		g.mainmenu.Draw(screen)

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
