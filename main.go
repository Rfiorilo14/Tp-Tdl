package main

import (
	"snake-game/game"
)

func main() {
	// Crear un tablero de juego de 20x10
	board := game.NewBoard(80, 20)

	// Dibujar el tablero en la consola
	board.Display() // Cambiado a Display() en lugar de Draw()
}
