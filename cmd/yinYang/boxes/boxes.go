package boxes

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const BoxDefaultSize = 20

type Box struct {
	X     float32
	Y     float32
	Color color.RGBA
}

func (b *Box) Draw(screen *ebiten.Image) {
	vector.DrawFilledRect(screen, b.X, b.Y, BoxDefaultSize, BoxDefaultSize, b.Color, false)
}

func New(x float32, y float32, color color.RGBA) Box {
	return Box{
		X:     x,
		Y:     y,
		Color: color,
	}
}
