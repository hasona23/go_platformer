package ui

import (
	"bytes"
	"go_platformer/components"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type Label struct {
	Text  string
	Pos   [2]int
	Style Style
}
type Style struct {
	Font        *text.GoTextFaceSource
	Size        int
	Color       color.Color
	LineSpacing uint32
}

func NewLabel(txt string, x, y int, fontFile []byte, size int, color color.Color) *Label {
	font, err := text.NewGoTextFaceSource(bytes.NewReader(fontFile))
	if err != nil {
		log.Fatalf("Error loading font file %v", err)
	}

	return &Label{
		Text:  txt,
		Pos:   [2]int{x, y},
		Style: NewStyle(font, size, color),
	}
}
func NewStyle(font *text.GoTextFaceSource, size int, color color.Color) Style {
	return Style{
		Font:        font,
		Size:        size,
		Color:       color,
		LineSpacing: 1,
	}
}
func (s Style) Face() *text.GoTextFace {
	return &text.GoTextFace{Source: s.Font, Size: float64(s.Size)}
}
func (l *Label) SetText(txt string) {
	l.Text = txt
}
func (l *Label) SetStyle(style Style) {
	l.Style = style
}
func (l *Label) SetPosition(x, y int) {
	l.Pos = [2]int{x, y}
}
func (l *Label) GetText() string {
	return l.Text
}

func (l *Label) GetPosition() [2]int {
	return l.Pos
}

func (l *Label) GetStyle() Style {
	return l.Style
}
func (l *Label) Draw(screen *ebiten.Image) {
	op := &text.DrawOptions{}
	op.GeoM.Translate(float64(l.Pos[0]), float64(l.Pos[1]))
	op.ColorScale.ScaleWithColor(l.Style.Color)
	text.Draw(screen, l.Text, l.Style.Face(), op)
}
func (l *Label) DrawCam(screen *ebiten.Image, cam components.Camera) {
	op := &text.DrawOptions{}
	op.GeoM.Translate(float64(l.Pos[0]+cam.X), float64(l.Pos[1])+float64(cam.Y))
	op.ColorScale.ScaleWithColor(l.Style.Color)
	text.Draw(screen, l.Text, l.Style.Face(), op)
}
func (l *Label) CenterText() {
	w, h := l.GetDimensions()
	l.Pos[0] -= int(w / 2)
	l.Pos[1] -= int(h / 2)
}
func (l *Label) GetDimensions() (width, height int) {
	f := l.Style.Face()
	w, h := text.Measure(l.Text, f, float64(l.Style.LineSpacing))
	return int(w), int(h)
}
func (l *Label) SetColor(c color.Color) {
	l.Style.Color = c
}
func (l *Label) GetBounds() (x, y, width, height int) {
	width, height = l.GetDimensions()
	return l.Pos[0], l.Pos[1], width, height
}
func (l *Label) SetFontSize(size int) {
	l.Style.Size = size
}
func (l *Label) Move(dx, dy int) {
	l.Pos[0] += dx
	l.Pos[1] += dy
}

func (l *Label) MoveX(dx int) {
	l.Pos[0] += dx
}

func (l *Label) MoveY(dy int) {
	l.Pos[1] += dy
}
