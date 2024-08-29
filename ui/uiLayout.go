package ui

import (
	"go_platformer/components"
	"slices"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type UILayout struct {
	id            string
	labels        map[string]*Label
	buttons       map[string]*Button
	focusedButton string
}

func NewUILayout(id string) *UILayout {
	return &UILayout{
		id:      id,
		labels:  make(map[string]*Label),
		buttons: make(map[string]*Button),
	}
}

func (u *UILayout) GetID() string {
	return u.id
}
func (u *UILayout) AddButton(name string, button *Button) {
	u.buttons[name] = button
	if u.focusedButton == "" {
		u.focusedButton = name
	}
}
func (u *UILayout) AddLabel(name string, label *Label) {
	u.labels[name] = label
}
func (u *UILayout) RemoveButton(name string) {
	delete(u.buttons, name)
}
func (u *UILayout) RemoveLabel(name string) {
	delete(u.labels, name)
}
func (u UILayout) GetButton(name string) (*Button, bool) {
	return u.buttons[name], u.HasButton(name)
}
func (u UILayout) GetLabel(name string) (*Label, bool) {
	return u.labels[name], u.HasLabel(name)
}
func (u *UILayout) GetButtonNames() []string {
	keys := make([]string, 0, len(u.buttons))
	for k := range u.buttons {
		keys = append(keys, k)
	}
	slices.Sort(keys)
	return keys
}

func (u *UILayout) GetLabelNames() []string {
	keys := make([]string, 0, len(u.labels))
	for k := range u.labels {
		keys = append(keys, k)
	}
	slices.Sort(keys)
	return keys
}
func (u *UILayout) HasButton(name string) bool {
	_, exists := u.buttons[name]
	return exists
}

func (u *UILayout) HasLabel(name string) bool {
	_, exists := u.labels[name]
	return exists
}
func (u *UILayout) Clear() {
	u.buttons = make(map[string]*Button)
	u.labels = make(map[string]*Label)
	u.focusedButton = ""
}
func (u *UILayout) GetFocusedButton() (*Button, bool) {
	return u.GetButton(u.focusedButton)
}

// checks for buttons hover/pressing and calls actions responding to this
// also allow navigation through keys by arrows or wasd
func (u *UILayout) Update() {
	u.navigation()
	u.updateButtons()
}

// function called in update method and allows you to navigate between button by arrow keys or wasd
func (u *UILayout) navigation() {

	keys := make([]string, 0, len(u.buttons))
	for k := range u.buttons {
		keys = append(keys, k)
	}
	if len(keys) == 0 {
		return
	}

	if u.focusedButton == "" {
		u.focusedButton = keys[0]
	}

	index := slices.Index(keys, u.focusedButton)

	for i, b := range u.buttons {
		if b.IsHover() {
			index = slices.Index(keys, i)
			break
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) || inpututil.IsKeyJustPressed(ebiten.KeyW) {
		index = (index - 1 + len(keys)) % len(keys)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) || inpututil.IsKeyJustPressed(ebiten.KeyS) {
		index = (index + 1) % len(keys)
	}

	u.focusedButton = keys[index]
}

func (u *UILayout) updateButtons() {
	for _, b := range u.buttons {
		if b == u.buttons[u.focusedButton] {
			b.OnHover()
			if b.IsPressed() || inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
				b.OnClick()
			}
		} else {
			b.DefaultColor()
		}
	}
}
func (u *UILayout) Draw(screen *ebiten.Image) {
	for _, button := range u.buttons {
		button.Draw(screen)
	}
	for _, label := range u.labels {
		label.Draw(screen)
	}
}
func (u *UILayout) DrawCam(screen *ebiten.Image, cam components.Camera) {
	for _, button := range u.buttons {
		button.DrawCam(screen, cam)
	}
	for _, label := range u.labels {
		label.DrawCam(screen, cam)
	}
}
func (u *UILayout) ApplyHoverToAllButtons(hoverEffect func(b *Button)) {
	for i := range u.buttons {
		u.buttons[i].AddHoverEvent(hoverEffect)
	}
}
func (u *UILayout) ApplyClickToAllButtons(hoverEffect func(b *Button)) {
	for i := range u.buttons {
		u.buttons[i].AddClickEvent(hoverEffect)
	}
}
