package ui

import (
	"go_platformer/components"
	"slices"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type UIManager struct {
	labels        []*Label
	buttons       []*Button
	focusedButton int
}

func NewUIManager() *UIManager {
	return &UIManager{}
}
func (u *UIManager) AddButton(button *Button) {
	u.buttons = append(u.buttons, button)
}
func (u *UIManager) AddLabel(label *Label) {
	u.labels = append(u.labels, label)
}

// checks for buttons hover/pressing and calls actions responding to this
// also allow navigation through keys by arrows or wasd
func (u *UIManager) Update() {
	u.navigation()
	u.updateHoveredButtons()
	u.updatePressedButtons()
}

// function called in update method and allows you to navigate between button by arrow keys or wasd
func (u *UIManager) navigation() {
	if (inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) || inpututil.IsKeyJustPressed(ebiten.KeyW)) && u.focusedButton > 0 {
		u.focusedButton--
	}
	if (inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) || inpututil.IsKeyJustPressed(ebiten.KeyS)) && u.focusedButton < len(u.buttons)-1 {
		u.focusedButton++
	}

}
func (u *UIManager) updateHoveredButtons() {
	for _, b := range u.buttons {
		if b.IsHover() {
			u.focusedButton = slices.Index(u.buttons, b)
		}
		if b == u.buttons[u.focusedButton] {
		}
	}
}
func (u *UIManager) updatePressedButtons() {
	for _, b := range u.buttons {
		if b == u.buttons[u.focusedButton] && (b.IsPressed() || inpututil.IsKeyJustPressed(ebiten.KeyEnter)) {

		}
	}
}
func (u *UIManager) Draw(screen *ebiten.Image) {
	for _, button := range u.buttons {
		button.Draw(screen)
	}
	for _, label := range u.labels {
		label.Draw(screen)
	}
}
func (u *UIManager) DrawCam(screen *ebiten.Image, cam components.Camera) {
	for _, button := range u.buttons {
		button.DrawCam(screen, cam)
	}
	for _, label := range u.labels {
		label.DrawCam(screen, cam)
	}
}
