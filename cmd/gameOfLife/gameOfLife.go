package main

import (
	"math/rand"
	"time"

	"github.com/charmbracelet/log"
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	aliveColor   = 0xff
	deadColor    = 0x0f
	screenWidth  = 400
	screenHeight = 400
)

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

type Grid struct {
	width  int
	height int
	pixels []byte
	cells  []bool
	rng    *rand.Rand
}

func NewGrid(width, height int, density float32) Grid {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cells := make([]bool, width*height)
	for i := range cells {
		cells[i] = rng.Float32() < density
	}

	newGrid := Grid{
		width:  width,
		height: height,
		pixels: make([]byte, width*height*4),
		cells:  cells,
		rng:    rng,
	}
	return newGrid
}

func (game *Grid) countNeighbour(x int, y int) int {
	// If is alive set to -1 to avoid skiping over the cell
	neighbours := boolToInt(game.cells[y*game.width+x]) * -1
	for yOffset := -1; yOffset <= 1; yOffset++ {
		for xOffset := -1; xOffset <= 1; xOffset++ {
			x2 := (x + xOffset + game.width) % game.width
			y2 := (y + yOffset + game.height) % game.height
			if game.cells[y2*game.width+x2] {
				neighbours++
			}
		}
	}
	return neighbours
}

func (game *Grid) updatePixelColor(status bool, position int) {
	pixOffset := 4 * position
	if status {
		game.pixels[pixOffset] = aliveColor
		game.pixels[pixOffset+1] = aliveColor
		game.pixels[pixOffset+2] = aliveColor
	} else {
		game.pixels[pixOffset] = deadColor
		game.pixels[pixOffset+1] = deadColor
		game.pixels[pixOffset+2] = deadColor
	}
	// Set the alpha chanel to MAX
	game.pixels[pixOffset+3] = 0xff
}

func (game *Grid) Update() {
	next := make([]bool, game.width*game.height)
	for y := 0; y < game.height; y++ {
		for x := 0; x < game.width; x++ {
			neighbours := game.countNeighbour(x, y)
			position := y*game.width + x
			if neighbours < 2 || 3 < neighbours {
				next[position] = false
				game.updatePixelColor(false, position)
			} else if (neighbours == 2 && game.cells[position]) || neighbours == 3 {
				next[position] = true
				game.updatePixelColor(true, position)
			}
		}
	}
	game.cells = next
}

func (game *Grid) CreateCells(x int, y int, density float32, area int) {
	if x < 0 || game.width < x || y < 0 || game.height < y {
		return
	}
	for yOffset := -area; yOffset <= area; yOffset++ {
		for xOffset := -area; xOffset <= area; xOffset++ {
			x2 := (x + xOffset + game.width) % game.width
			y2 := (y + yOffset + game.height) % game.height
			status := game.rng.Float32() < density
			position := y2*game.width + x2
			game.cells[position] = status
			game.updatePixelColor(status, position)
		}
	}
}

type Game struct {
	grid Grid
}

func (game *Game) Update() error {
	game.grid.Update()
	if ebiten.IsMouseButtonPressed(ebiten.MouseButton0) {
		x, y := ebiten.CursorPosition()
		game.grid.CreateCells(x, y, 0.20, 5)
	}
	return nil
}

func (game *Game) Draw(screen *ebiten.Image) {
	screen.WritePixels(game.grid.pixels)
}

func (game *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetWindowTitle("Game of Life")

	g := Game{
		grid: NewGrid(screenWidth, screenHeight, 0.10),
	}

	if err := ebiten.RunGame(&g); err != nil {
		log.Fatal(err)
	}
}
