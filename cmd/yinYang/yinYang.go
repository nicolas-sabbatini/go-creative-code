package main

import (
	"github.com/charmbracelet/log"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/nicolas-sabbatini/go-creative-code/cmd/yinYang/balls"
	"github.com/nicolas-sabbatini/go-creative-code/cmd/yinYang/boxes"
	"github.com/nicolas-sabbatini/go-creative-code/cmd/yinYang/globals"
	"github.com/nicolas-sabbatini/go-creative-code/cmd/yinYang/physics"
)

type Game struct {
	boxes    []boxes.Box
	balls    []balls.Ball
	gridSize int
}

func newGame() Game {
	newBoxes := make([]boxes.Box, 0)
	for x := 0; x < globals.ScreenWidth; x += boxes.BoxDefaultSize {
		for y := 0; y < globals.ScreenHeight; y += boxes.BoxDefaultSize {
			if x < globals.ScreenWidth/2 {
				newBoxes = append(newBoxes, boxes.New(float32(x), float32(y), globals.Black))
			} else {
				newBoxes = append(newBoxes, boxes.New(float32(x), float32(y), globals.White))
			}
		}
	}
	newBalls := make([]balls.Ball, 0)
	newBalls = append(newBalls, balls.New(float32(globals.ScreenWidth)*0.25, float32(globals.ScreenHeight)*0.5, globals.White, globals.Black, 1.0))
	newBalls = append(newBalls, balls.New(float32(globals.ScreenWidth)*0.75, float32(globals.ScreenHeight)*0.5, globals.Black, globals.White, -1.0))
	if globals.ScreenWidth/boxes.BoxDefaultSize != globals.ScreenHeight/boxes.BoxDefaultSize {
		log.Fatal("The sizes are not the same")
	}
	return Game{
		boxes:    newBoxes,
		balls:    newBalls,
		gridSize: globals.ScreenWidth / boxes.BoxDefaultSize,
	}
}

func (game *Game) Update() error {
	for i := 0; i < len(game.balls); i++ {
		game.balls[i].Update()
		physics.Colide(&game.balls[i], &game.boxes, game.gridSize)
	}
	return nil
}

func (game *Game) Draw(screen *ebiten.Image) {
	screen.Clear()
	for i := 0; i < len(game.boxes); i++ {
		game.boxes[i].Draw(screen)
	}
	for i := 0; i < len(game.balls); i++ {
		game.balls[i].Draw(screen)
	}
}

func (game *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return globals.ScreenWidth, globals.ScreenHeight
}

func main() {
	ebiten.SetWindowSize(globals.ScreenWidth, globals.ScreenHeight)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetWindowTitle("Yin Yang Pong")

	g := newGame()

	if err := ebiten.RunGame(&g); err != nil {
		log.Fatal(err)
	}
}
