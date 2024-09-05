package main

import (
	"go_platformer/assets"
	"go_platformer/components"
	"go_platformer/entities"
	"go_platformer/particles"
	"go_platformer/tilemap"
	ui "go_platformer/ui"
	"image"
	"image/color"
	"log"
	"os"
	"slices"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type GameState int

const (
	MainMenu GameState = iota
	Main
	Pause
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
	cam       components.Camera
	level1    *tilemap.Level
	enemies   []*entities.Enemy
	player    *entities.Player
	state     GameState
	mainmenu  *ui.UILayout
	gameUI    *ui.UILayout
	pauseUI   *ui.UILayout
	particles *particles.ParticleSystem
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

	ammoBar := ui.NewSpriteBar(assets.SpriteSheet.SubImage(image.Rect(3, 5*16+3, 14, 6*16-2)).(*ebiten.Image),
		10, 10, g.player.Ammo,
		components.Point{X: 5, Y: 5})

	g.gameUI = ui.NewUILayout("MainGameUI")
	g.gameUI.AddBar("ammo", ammoBar)
	start := ui.NewButton("START", 100, 100, 16, 2, assets.PixelFont, color.White, color.RGBA{255, 222, 206, 255}, color.White)
	save := ui.NewButton("SAVE", 100, 150, 16, 2, assets.PixelFont, color.White, color.RGBA{255, 222, 206, 255}, color.White)
	exit := ui.NewButton("QUIT", 100, 200, 16, 2, assets.PixelFont, color.White, color.RGBA{255, 222, 206, 255}, color.White)
	start.Style.BorderThickness = 3
	save.Style.BorderThickness = 3
	exit.Style.BorderThickness = 3
	g.mainmenu = ui.NewUILayout("MainMenu")
	g.mainmenu.AddButton(Start, start)
	g.mainmenu.AddButton(Save, save)
	g.mainmenu.AddButton(Exit, exit)
	g.mainmenu.ApplyHoverToAllButtons(hover)
	start.AddClickEvent(func(b *ui.Button) {
		if g.state != Main {
			g.state = Main
		}
	})

	exit.AddClickEvent(func(b *ui.Button) {
		os.Exit(0)
	})

	pauseText := ui.NewLabel("PAUSED!", DisplayWidth/2, DisplayHeight/2, assets.PixelFont, 32, color.White)
	pauseText.CenterText()
	g.pauseUI = ui.NewUILayout("pauseUI")
	g.pauseUI.AddLabel("pauseText", pauseText)
	returnButton := ui.NewButton("Back to Menu", 160, 175, 16, 2, assets.PixelFont, color.White, color.Transparent, color.White)
	returnButton.Style.TextOrientation = ui.Middle

	returnButton.AddHoverEvent(hover)
	returnButton.AddClickEvent(func(b *ui.Button) {
		g.state = MainMenu
	})
	returnButton.Centre()
	g.pauseUI.AddButton("tomenuButton", returnButton)
	g.particles = particles.NewParticleSystem(0, 0, g.level1.GetSizeInPixels()[0], g.level1.GetSizeInPixels()[1], 10, particles.RandomDirections,
		*particles.NewParticle(0, 0, particles.WithColor(color.RGBA{180, 211, 75, 255}), particles.WithScale(4), particles.WithSpeed(0.2)))
	g.particles.Decelration = 0.001
	g.particles.ParticleSpawnCount = 10
	g.particles.SpawnTime = components.NewTimer(60)
	g.particles.IsLooped = true

}

func hover(b *ui.Button) {
	b.Text.Style.Color = color.Black
	b.Style.BorderColor = color.White
}

func (g *Game) Update() error {
	if g.state == Pause {
		g.pauseUI.Update()
	} else if g.state == MainMenu {
		g.mainmenu.Update()

	}
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		if g.state == Main {
			g.state = Pause
		} else if g.state == Pause {
			g.state = Main
		}
	}

	if g.state == Main {
		if b, _ := g.mainmenu.GetButton(Start); b.IsPressed() {
			return nil
		}
		if g.player.Died {
			g.enemies = []*entities.Enemy{}
			g.Init()
			g.state = Main
		}
		bar, _ := g.gameUI.GetBar("ammo")
		bar.SetValue(g.player.Ammo)
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

		g.enemies = slices.DeleteFunc(g.enemies, func(e *entities.Enemy) bool { return e.Dead })
		g.particles.Update()
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	if g.state != MainMenu {
		g.level1.DrawCamera(screen, assets.SpriteSheet, g.cam, false)
		for _, object := range g.enemies {
			object.Draw(screen, g.cam)
		}
		g.player.Draw(screen, g.cam)
		g.player.PhysicsEntity.Draw(screen, g.cam)
		g.gameUI.Draw(screen)
		g.particles.DrawCam(screen, g.cam)
	}
	if g.state == Pause {
		vector.DrawFilledRect(screen, 0, 0, DisplayWidth, DisplayHeight, color.RGBA{60, 60, 60, 100}, false)
		g.pauseUI.Draw(screen)
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
