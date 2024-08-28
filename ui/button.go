package ui

import (
	"go_platformer/components"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// text orientation in button
type Orientation int

const (
	TopLeft Orientation = iota
	Middle
	BottomRight
)

type Button struct {
	sprite  *ebiten.Image
	Text    Label
	rect    components.Rect
	OnClick func(b *Button)
	OnHover func(b *Button)
	ButtonStyle
}

func NewSpriteButton(sprite *ebiten.Image, text string, x, y, fontSize, scale int, fontFile []byte, textColor color.Color) *Button {

	button := &Button{}
	button.Text = *NewLabel(text, x, y, fontFile, fontSize, color.Black)
	button.ButtonStyle.Color = color.Transparent
	button.sprite = sprite
	button.rect = components.NewRect(x, y, scale*sprite.Bounds().Dx(), scale*sprite.Bounds().Dy())
	button.Text.Style.Color = textColor
	button.defaultTextColor = textColor
	button.OnClick = func(b *Button) {}
	button.OnHover = func(b *Button) {}
	return button
}

// normal square/rectangle button with border background color
func NewButton(txt string, x, y, fontSize int, scale float64, fontFile []byte, textColor, backColor, bordercolor color.Color) *Button {

	button := &Button{}
	button.Text = *NewLabel(txt, x, y, fontFile, fontSize, color.Black)
	button.ButtonStyle.X = x
	button.ButtonStyle.Y = y
	button.ButtonStyle.Scale = scale
	button.ButtonStyle.BorderThickness = 1
	button.ButtonStyle.Color = color.Transparent
	button.ButtonStyle.BorderColor = bordercolor
	button.ButtonStyle.BackColor = backColor
	button.ButtonStyle.defaultBackColor = backColor
	button.ButtonStyle.defaultBorderColor = bordercolor
	button.Text.Style.Color = textColor
	button.defaultTextColor = textColor
	button.OnClick = func(b *Button) {}
	button.OnHover = func(b *Button) {}
	return button
}

// Button ButtonStyle is made Primarily for when there is no sprite or make effects for sprite
type ButtonStyle struct {
	X, Y               int
	Color              color.Color
	BorderColor        color.Color
	BackColor          color.Color
	defaultBackColor   color.Color
	defaultBorderColor color.Color
	defaultTextColor   color.Color
	Scale              float64
	BorderThickness    int
	TextOrientation    Orientation
}

// Draw button
func (b *Button) Draw(screen *ebiten.Image) {
	//check if has sprite otherwise will draw border(normal button)
	b.drawButton(screen)
	b.drawButtonText(screen)
}

func (b *Button) drawButton(screen *ebiten.Image) {
	if b.sprite != nil {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(b.rect.X), float64(b.rect.Y))
		op.GeoM.Scale(float64(b.ButtonStyle.Scale), float64(b.ButtonStyle.Scale))
		op.ColorScale.ScaleWithColor(b.ButtonStyle.Color)
		screen.DrawImage(b.sprite, op)
	} else {
		vector.DrawFilledRect(screen, float32(b.rect.X), float32(b.rect.Y), float32(b.rect.Width), float32(b.rect.Height), b.ButtonStyle.BackColor, false)
		vector.StrokeRect(screen, float32(b.rect.X), float32(b.rect.Y), float32((b.rect.Width + b.BorderThickness/2)), float32(b.rect.Height+b.BorderThickness/2),
			float32(b.ButtonStyle.BorderThickness),
			b.ButtonStyle.BorderColor, false)
	}
}
func (b *Button) drawButtonText(screen *ebiten.Image) {
	opText := &text.DrawOptions{}
	f := &text.GoTextFace{Source: b.Text.Style.Font, Size: float64(b.Text.Style.Size)}
	width, height := text.Measure(b.Text.Text, f, 1)

	opText.GeoM.Translate(float64(b.rect.X), float64(b.rect.Y))
	if b.Scale != 1 {
		switch b.ButtonStyle.TextOrientation {
		case Middle:
			opText.GeoM.Translate(float64(b.rect.Width)/2-width/2, float64(b.rect.Height)/2-height/2)

		case BottomRight:
			opText.GeoM.Translate(float64(b.rect.Width)-width*1.25, float64(b.rect.Height)-height*1.25)
			//opText.GeoM.Translate(float64(b.rect.Width)-width, float64(b.rect.Height)-height)

		default:
			//top left
			opText.GeoM.Translate(width*.25, height*.25)

		}
	}
	opText.ColorScale.ScaleWithColor(b.Text.Style.Color)
	text.Draw(screen, b.Text.Text, &text.GoTextFace{Source: b.Text.Style.Font, Size: float64(b.Text.Style.Size)}, opText)
	b.rect = components.NewRect(b.ButtonStyle.X, b.ButtonStyle.Y, int(width*b.ButtonStyle.Scale), int(height*b.ButtonStyle.Scale))

}

// draw button with camera in main so not fixed to screen
func (b *Button) DrawCam(screen *ebiten.Image, cam components.Camera) {
	//check if has sprite otherwise will draw border(normal button)
	b.drawButtonCam(screen, cam)
	b.drawButtonTextCam(screen, cam)
}
func (b *Button) drawButtonCam(screen *ebiten.Image, cam components.Camera) {
	if b.sprite != nil {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(b.rect.X), float64(b.rect.Y))
		op.GeoM.Translate(float64(cam.X), float64(cam.Y))
		op.GeoM.Scale(float64(b.ButtonStyle.Scale), float64(b.ButtonStyle.Scale))
		op.ColorScale.ScaleWithColor(b.ButtonStyle.Color)
		screen.DrawImage(b.sprite, op)
	} else {
		vector.DrawFilledRect(screen, float32(b.rect.X+cam.X), float32(b.rect.Y+cam.Y), float32(b.rect.Width), float32(b.rect.Height), b.ButtonStyle.BackColor, false)
		vector.StrokeRect(screen, float32(b.rect.X+cam.X), float32(b.rect.Y+cam.Y), float32((b.rect.Width + b.BorderThickness/2)), float32(b.rect.Height+b.BorderThickness/2),
			float32(b.ButtonStyle.BorderThickness),
			b.ButtonStyle.BorderColor, false)
	}
}
func (b *Button) drawButtonTextCam(screen *ebiten.Image, cam components.Camera) {
	opText := &text.DrawOptions{}
	f := &text.GoTextFace{Source: b.Text.Style.Font, Size: float64(b.Text.Style.Size)}
	width, height := text.Measure(b.Text.Text, f, 1)

	opText.GeoM.Translate(float64(b.rect.X), float64(b.rect.Y))
	opText.GeoM.Translate(float64(cam.X), float64(cam.Y))
	if b.Scale != 1 {
		switch b.ButtonStyle.TextOrientation {
		case Middle:
			opText.GeoM.Translate(float64(b.rect.Width)/2-width/2, float64(b.rect.Height)/2-height/2)

		case BottomRight:
			opText.GeoM.Translate(float64(b.rect.Width)-width*1.25, float64(b.rect.Height)-height*1.25)
			//opText.GeoM.Translate(float64(b.rect.Width)-width, float64(b.rect.Height)-height)

		default:
			//top left
			opText.GeoM.Translate(width*.25, height*.25)

		}
	}
	opText.ColorScale.ScaleWithColor(b.Text.Style.Color)
	text.Draw(screen, b.Text.Text, &text.GoTextFace{Source: b.Text.Style.Font, Size: float64(b.Text.Style.Size)}, opText)
	b.rect = components.NewRect(b.ButtonStyle.X, b.ButtonStyle.Y, int(width*b.ButtonStyle.Scale), int(height*b.ButtonStyle.Scale))

}

// check if button is being hovered on by the mouse cursor
func (b *Button) IsHover() bool {
	return b.rect.Contains(ebiten.CursorPosition())
}

// check if button is pressed by mouse
func (b *Button) IsPressed() bool {
	return (b.IsHover() && inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft))
}
func (b *Button) DefaultColor() {
	b.ButtonStyle.Color = color.Transparent
	b.ButtonStyle.BorderColor = b.ButtonStyle.defaultBorderColor
	b.ButtonStyle.BackColor = b.ButtonStyle.defaultBackColor
	b.Text.Style.Color = b.defaultTextColor
}
