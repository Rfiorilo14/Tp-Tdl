package game

import (
	"fmt"
)

// Constantes para representar diferentes elementos en el tablero
const (
	EmptyCell    = " " // Celda vacía
	SnakeCell    = "S" // Representación de la serpiente
	FoodCell     = "F" // Representación de la comida
	ObstacleCell = "O" // Representación de un obstáculo
	BorderCell   = "-" // Representación del borde
)

// Struct para el tablero
type Board struct {
	Width, Height int
	Grid          [][]string // Matriz 2D que representa el tablero
}

// NewBoard crea un nuevo tablero con las dimensiones especificadas
func NewBoard(width, height int) *Board {
	board := &Board{
		Width:  width,
		Height: height,
		Grid:   make([][]string, height),
	}

	// Inicializa el tablero con celdas vacías
	board.Reset()
	return board
}

// Reset limpia el tablero, dejando solo celdas vacías
func (b *Board) Reset() {
	for i := range b.Grid {
		b.Grid[i] = make([]string, b.Width)
		for j := range b.Grid[i] {
			b.Grid[i][j] = EmptyCell
		}
	}
}

// PlaceSnake coloca la serpiente en el tablero en una posición específica
func (b *Board) PlaceSnake(x, y int) {
	if b.isInBounds(x, y) {
		b.Grid[y][x] = SnakeCell
	}
}

// PlaceFood coloca comida en el tablero en una posición específica
func (b *Board) PlaceFood(x, y int) {
	if b.isInBounds(x, y) {
		b.Grid[y][x] = FoodCell
	}
}

// PlaceObstacle coloca un obstáculo en el tablero en una posición específica
func (b *Board) PlaceObstacle(x, y int) {
	if b.isInBounds(x, y) {
		b.Grid[y][x] = ObstacleCell
	}
}

// isInBounds verifica si una posición está dentro de los límites del tablero
func (b *Board) isInBounds(x, y int) bool {
	return x >= 0 && x < b.Width && y >= 0 && y < b.Height
}

// Display muestra el tablero en la terminal, con bordes
func (b *Board) Display() {
	// Imprimir el borde superior
	fmt.Print(BorderCell)
	for i := 0; i < b.Width; i++ {
		fmt.Print(BorderCell)
	}
	fmt.Println(BorderCell)

	// Imprimir las filas del tablero con bordes laterales
	for _, row := range b.Grid {
		fmt.Print(BorderCell) // Borde izquierdo
		for _, cell := range row {
			fmt.Print(cell)
		}
		fmt.Println(BorderCell) // Borde derecho
	}

	// Imprimir el borde inferior
	fmt.Print(BorderCell)
	for i := 0; i < b.Width; i++ {
		fmt.Print(BorderCell)
	}
	fmt.Println(BorderCell)
}
