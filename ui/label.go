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
	Font  *text.GoTextFaceSource
	Size  int
	Color color.Color
}

func NewLabel(txt string, x, y int, fontFile []byte, size int, color color.Color) *Label {
	label := &Label{}
	label.Text = txt
	label.Pos = [2]int{x, y}
	font, err := text.NewGoTextFaceSource(bytes.NewReader(fontFile))
	if err != nil {
		log.Fatal("Error loading font file :", err)
	}
	label.Style.Font = font
	label.Style.Size = size
	label.Style.Color = color
	return label
}

func (l *Label) Draw(screen *ebiten.Image) {
	op := &text.DrawOptions{}
	op.GeoM.Translate(float64(l.Pos[0]), float64(l.Pos[1]))
	op.ColorScale.ScaleWithColor(l.Style.Color)
	text.Draw(screen, l.Text, &text.GoTextFace{Source: l.Style.Font, Size: float64(l.Style.Size)}, op)
}
func (l *Label) DrawCam(screen *ebiten.Image, cam components.Camera) {
	op := &text.DrawOptions{}
	op.GeoM.Translate(float64(l.Pos[0]+cam.X), float64(l.Pos[1])+float64(cam.Y))

	text.Draw(screen, l.Text, &text.GoTextFace{Source: l.Style.Font, Size: float64(l.Style.Size)}, op)
}
