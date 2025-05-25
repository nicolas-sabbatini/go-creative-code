package balls

import (
	"image/color"
	"math/rand/v2"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/nicolas-sabbatini/go-creative-code/cmd/yinYang/globals"
)

const BallSize = 10

type Speed struct {
	X float32
	Y float32
}

type Ball struct {
	X      float32
	Y      float32
	Color  color.RGBA
	Change color.RGBA
	Speed  Speed
}

func (b *Ball) Draw(screen *ebiten.Image) {
	vector.DrawFilledCircle(screen, b.X, b.Y, BallSize, b.Color, false)
}

func (b *Ball) Update() {
	b.X += b.Speed.X
	b.Y += b.Speed.Y
	if b.X-BallSize < 0 {
		b.X = 0 + BallSize
		b.MirrorX()
	}
	if b.X+BallSize > globals.ScreenWidth {
		b.X = globals.ScreenWidth - BallSize
		b.MirrorX()
	}
	if b.Y-BallSize < 0 {
		b.Y = 0 + BallSize
		b.MirrorY()
	}
	if b.Y+BallSize > globals.ScreenWidth {
		b.Y = globals.ScreenWidth - BallSize
		b.MirrorY()
	}
}

func (b *Ball) MirrorX() {
	b.Speed.X = -b.Speed.X + rand.Float32() - 0.5
}

func (b *Ball) MirrorY() {
	b.Speed.Y = -b.Speed.Y + rand.Float32() - 0.5
}

func New(x float32, y float32, color color.RGBA, change color.RGBA, speedDir float32) Ball {
	return Ball{
		X:      x,
		Y:      y,
		Color:  color,
		Change: change,
		Speed: Speed{
			X: 6 * speedDir,
			Y: 6 * speedDir,
		},
	}
}
