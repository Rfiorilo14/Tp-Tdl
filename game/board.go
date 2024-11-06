// game/board.go
package game

import (
	"math/rand"
	"time"
)

const (
	EmptyCell    = " " // Celda vacía
	SnakeCell    = "S" // Representación de la serpiente
	FoodCell     = "F" // Representación de la comida
	ObstacleCell = "O" // Representación de un obstáculo
	BorderCell   = "-" // Representación del borde
)

type Board struct {
	Width, Height int
	Grid          [][]string // Matriz 2D que representa el tablero
}

func NewBoard(width, height int) *Board {
	board := &Board{
		Width:  width,
		Height: height,
		Grid:   make([][]string, height),
	}
	board.Reset()
	board.PlaceRandomFood() // Colocar comida al iniciar
	return board
}

func (b *Board) Reset() {
	for i := range b.Grid {
		b.Grid[i] = make([]string, b.Width)
		for j := range b.Grid[i] {
			b.Grid[i][j] = EmptyCell
		}
	}
}

// PlaceRandomFood coloca comida en una posición aleatoria en el tablero
func (b *Board) PlaceRandomFood() {
	rand.Seed(time.Now().UnixNano())
	for {
		x, y := rand.Intn(b.Width), rand.Intn(b.Height)
		if b.Grid[y][x] == EmptyCell {
			b.PlaceFood(x, y)
			break
		}
	}
}

func (b *Board) PlaceSnake(x, y int) {
	if b.isInBounds(x, y) {
		b.Grid[y][x] = SnakeCell
	}
}

func (b *Board) PlaceFood(x, y int) {
	if b.isInBounds(x, y) {
		b.Grid[y][x] = FoodCell
	}
}

func (b *Board) PlaceObstacle(x, y int) {
	if b.isInBounds(x, y) {
		b.Grid[y][x] = ObstacleCell
	}
}

func (b *Board) isInBounds(x, y int) bool {
	return x >= 0 && x < b.Width && y >= 0 && y < b.Height
}
