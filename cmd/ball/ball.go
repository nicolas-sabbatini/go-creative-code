package main

import (
	_ "embed"

	"github.com/charmbracelet/log"
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	screenWidth  = 800
	screenHeight = 600
)

//go:embed assets/ball.kage
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
	opts.Uniforms = make(map[string]interface{})
	opts.Uniforms["Center"] = []float32{
		float32(screenWidth / 2),
		float32(screenHeight / 2),
	}
	opts.Images[0] = game.input
	game.ouput.DrawRectShader(screenWidth, screenHeight, game.shader, opts)

	screen.Clear()
	screen.DrawImage(game.ouput, nil)
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
