package main

import (
	_ "embed"
	"math/rand"
	"time"

	"github.com/charmbracelet/log"
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	screenWidth  = 1280 * 2
	screenHeight = 720 * 2
	lifeDensity  = 0.1
)

//go:embed assets/gameOfLife.kage
var shaderProgram []byte

type Game struct {
	shader *ebiten.Shader
	ouput  *ebiten.Image
	input  *ebiten.Image
}

func newGame() Game {
	shader, err := ebiten.NewShader(shaderProgram)
	if err != nil {
		log.Fatal(err)
	}
	input := ebiten.NewImage(screenWidth, screenHeight)
	ouput := ebiten.NewImage(screenWidth, screenHeight)

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cells := make([]byte, screenWidth*screenHeight*4)
	for x := 0; x < screenWidth; x++ {
		for y := 0; y < screenHeight; y++ {
			position := (y*screenWidth + x) * 4
			if rng.Float32() < lifeDensity {
				cells[position] = 0xff
				cells[position+1] = 0xff
				cells[position+2] = 0xff
			} else {
				cells[position] = 0x00
				cells[position+1] = 0x00
				cells[position+2] = 0x00
			}
			cells[position+3] = 0xff
		}
	}
	ouput.WritePixels(cells)

	return Game{
		shader: shader,
		ouput:  ouput,
		input:  input,
	}
}

func (game *Game) Update() error {
	temp := game.input
	game.input = game.ouput
	game.ouput = temp
	return nil
}

func (game *Game) Draw(screen *ebiten.Image) {
	opts := &ebiten.DrawRectShaderOptions{}
	opts.Images[0] = game.input
	game.ouput.DrawRectShader(screenWidth, screenHeight, game.shader, opts)

	screen.Clear()
	opts2 := &ebiten.DrawImageOptions{}
	opts2.Filter = ebiten.FilterNearest
	screen.DrawImage(game.ouput, opts2)
}

func (game *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetWindowTitle("Game of Life Shader")

	g := newGame()

	if err := ebiten.RunGame(&g); err != nil {
		log.Fatal(err)
	}
}
