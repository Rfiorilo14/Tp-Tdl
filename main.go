package main

import (
	"F:\tps_de_roy\FUIBA\TDL\snake-game\game"
)

func main() {
	// Crear un tablero de 10x10
	board := game.NewBoard(10, 10)

	// Mostrar el tablero en la terminal
	board.Display()
}
