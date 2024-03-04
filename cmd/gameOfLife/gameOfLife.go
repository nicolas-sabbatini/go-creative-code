package main

import (
	"log"
	"math/rand"
	"time"

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

func (self *Grid) countNeighbour(x int, y int) int {
	// If is alive set to -1 to avoid skiping over the cell
	neighbours := boolToInt(self.cells[y*self.width+x]) * -1
	for yOffset := -1; yOffset <= 1; yOffset++ {
		for xOffset := -1; xOffset <= 1; xOffset++ {
			x2 := (x + xOffset + self.width) % self.width
			y2 := (y + yOffset + self.height) % self.height
			if self.cells[y2*self.width+x2] {
				neighbours++
			}
		}
	}
	return neighbours
}

func (self *Grid) updatePixelColor(status bool, position int) {
	pixOffset := 4 * position
	if status {
		self.pixels[pixOffset] = aliveColor
		self.pixels[pixOffset+1] = aliveColor
		self.pixels[pixOffset+2] = aliveColor
	} else {
		self.pixels[pixOffset] = deadColor
		self.pixels[pixOffset+1] = deadColor
		self.pixels[pixOffset+2] = deadColor
	}
	// Set the alpha chanel to MAX
	self.pixels[pixOffset+3] = 0xff
}

func (self *Grid) Update() {
	next := make([]bool, self.width*self.height)
	for y := 0; y < self.height; y++ {
		for x := 0; x < self.width; x++ {
			neighbours := self.countNeighbour(x, y)
			position := y*self.width + x
			if neighbours < 2 || 3 < neighbours {
				next[position] = false
				self.updatePixelColor(false, position)
			} else if (neighbours == 2 && self.cells[position]) || neighbours == 3 {
				next[position] = true
				self.updatePixelColor(true, position)
			}
		}
	}
	self.cells = next
}

func (self *Grid) CreateCells(x int, y int, density float32, area int) {
	if x < 0 || self.width < x || y < 0 || self.height < y {
		return
	}
	for yOffset := -area; yOffset <= area; yOffset++ {
		for xOffset := -area; xOffset <= area; xOffset++ {
			x2 := (x + xOffset + self.width) % self.width
			y2 := (y + yOffset + self.height) % self.height
			status := self.rng.Float32() < density
			position := y2*self.width + x2
			self.cells[position] = status
			self.updatePixelColor(status, position)
		}
	}
}

type Game struct {
	grid Grid
}

func (self *Game) Update() error {
	self.grid.Update()
	if ebiten.IsMouseButtonPressed(ebiten.MouseButton0) {
		x, y := ebiten.CursorPosition()
		self.grid.CreateCells(x, y, 0.20, 5)
	}
	return nil
}

func (self *Game) Draw(screen *ebiten.Image) {
	screen.WritePixels(self.grid.pixels)
}

func (self *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
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
