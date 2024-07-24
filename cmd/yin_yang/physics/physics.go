package physics

import (
	"github.com/nicolas-sabbatini/project-name/cmd/yin_yang/balls"
	"github.com/nicolas-sabbatini/project-name/cmd/yin_yang/boxes"
)

func indexToGrid(i int, gridSize int) (int, int) {
	return i % gridSize, i / gridSize
}

func gridToIndex(x int, y int, gridSize int) int {
	x /= boxes.BoxDefaultSize
	y /= boxes.BoxDefaultSize
	return y + x*gridSize
}

func Colide(b *balls.Ball, boxes *[]boxes.Box, gridSize int) {
	points := [][]float32{{b.X, b.Y - balls.BallSize},
		{b.X, b.Y + balls.BallSize - 1},
		{b.X - balls.BallSize, b.Y},
		{b.X + balls.BallSize - 1, b.Y}}
	colitions := []bool{false, false, false, false}
	for c, point := range points {
		i := gridToIndex(int(point[0]), int(point[1]), gridSize)
		if (*boxes)[i].Color == b.Color {
			(*boxes)[i].Color = b.Change
			colitions[c] = true
		}
	}
	if colitions[2] || colitions[3] {
		b.MirrorX()
	}
	if colitions[0] || colitions[1] {
		b.MirrorY()
	}
}
