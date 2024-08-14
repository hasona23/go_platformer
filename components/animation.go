package components

import (
	"errors"

	"github.com/hajimehoshi/ebiten/v2"
)

type AnimSprite struct {
	Atlas         *ebiten.Image
	Img           *ebiten.Image
	width, height int
	Animations    map[string]AnimationFrame
	Current       AnimationFrame
	IsInverted    int
}

func (sprite AnimSprite) GetWidth() int {
	return sprite.Img.Bounds().Dx()
}

func (sprite AnimSprite) Origin(op *ebiten.DrawImageOptions) {
	s := sprite.Img.Bounds()
	op.GeoM.Translate(-float64(s.Dx())/2, -float64(s.Dy())/2)

}
func (sprite *AnimSprite) Rotate(op *ebiten.DrawImageOptions, angle float64) {
	s := sprite.Img.Bounds()
	op.GeoM.Translate(-float64(s.Dx())/2, -float64(s.Dy())/2)
	op.GeoM.Rotate(angle)
	op.GeoM.Translate(float64(s.Dx())/2, float64(s.Dy())/2)
}

func NewAnimeSprite(img *ebiten.Image, Width, Height int) AnimSprite {
	return AnimSprite{Atlas: img, Img: ebiten.NewImage(Width, Height), width: Width, height: Height, Animations: make(map[string]AnimationFrame)}
}

func (a *AnimSprite) Add(animtion AnimationFrame) {
	a.Animations[animtion.Name] = animtion
}

func (a *AnimSprite) playAnim() error {
	var err error = nil
	if a.Current.IsEmpty() {
		err = errors.New("animation is empty")
		return err
	}
	if _, ok := a.Animations[a.Current.Name]; ok {
		a.Current.Update()
	} else {
		err = errors.New("there is no animation with this name")
		return err
	}
	return nil
}
func (a *AnimSprite) ChangeAnim(name string) {
	a.Current = a.Animations[name]
}
func (a *AnimSprite) Update() error {
	err := a.playAnim()
	return err
}

// AnimeSprite Consist of many animation Frames to Player
type AnimationFrame struct {
	row_min     int
	row_max     int
	row_current int
	col_min     int
	col_max     int
	col_current int
	IsEnd       bool
	Name        string
	Timer
}

func (a AnimationFrame) IsEmpty() bool {
	return AnimationFrame{} != a
}

func NewAnimationFrame(row_min, row_max, col_min, col_max int, duration float64, name string) AnimationFrame {
	return AnimationFrame{row_min: row_min, row_max: row_max, row_current: row_min, col_min: col_min, col_max: col_max, col_current: col_min, IsEnd: false, Name: name, Timer: NewTimer(duration)}
}

// Updates Column and row number to affect the animsprite SubImage
func (a *AnimationFrame) Update() {
	a.Timer.UpdateTimer()
	//fmt.Println(a.Timer.GetCurrentTime())
	if a.Ticked() {

		a.IsEnd = false
		a.row_current++
		//fmt.Println(" Row", a.row_max, " ", a.row_current)
		if a.row_current >= a.row_max {

			a.row_current = a.row_min
			a.col_current++

		}
		//fmt.Println(" Col ", a.col_max, " ", a.col_current)
		if a.col_current >= a.col_max {

			a.col_current = a.col_min
			a.IsEnd = true

		}

	}
}
