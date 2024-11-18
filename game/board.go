// game/board.go
package game

import (
	"math/rand"
	pckg_snake "snake-game/snake"
	"time"
)

type Snake = pckg_snake.Snake

/*
Los elementos del tablero pueden ser numeros (en principio):

	0 : celda vacia
	1: celda ocupada por cuerpo snake de jugador 1
	2: celda ocupada por cuerpo snake de jugador 2
	3: celda ocupada por cuerpo snake de jugador 3
	4: un tipo de comida
	5: un tipo de obstaculo
	6: un tipo de power up

las comidas pueden implementar una interfaz general y cada tipo de comida puede tener diferentes
caracteristicas . Lo mismo para los obstaculos
*/
const (
	EmptyCell    = 0 // Celda vacía
	SnakeCell    = 1 // Representación de la serpiente
	FoodCell     = 2 // Representación de la comida
	ObstacleCell = 3 // Representación de un obstáculo
	BorderCell   = 4 // Representación del borde
)

type Board struct {
	Width, Height int
	Grid          [][]int // Matriz 2D que representa el tablero
}

func NewBoard(width, height int) *Board {
	board := &Board{
		Width:  width,
		Height: height,
		Grid:   make([][]int, height),
	}
	board.Reset()
	board.PlaceRandomFood() // Colocar comida al iniciar
	return board
}

func (b *Board) Reset() {
	for i := range b.Grid {
		b.Grid[i] = make([]int, b.Width)
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

func (b *Board) PlaceSnake(snake *Snake) {
	/*
		if b.isInBounds(x, y) {
			b.Grid[y][x] = SnakeCell
		}
	*/
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
