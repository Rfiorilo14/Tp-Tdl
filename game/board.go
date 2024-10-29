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
	BorderCell   = "#" // Representación del borde
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

	for i := range board.Grid {
		board.Grid[i] = make([]string, width)
		for j := range board.Grid[i] {
			board.Grid[i][j] = EmptyCell
		}
	}
	return board
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
