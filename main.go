// main.go
package main

import (
	"log"
	"snake-game/game"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	// Crear un tablero de juego de 120x60
	board := game.NewBoard(120, 60)
	gameInstance := game.NewGame(board)

	// Configuración de la ventana y ejecución del juego
	ebiten.SetWindowSize(1600, 800)
	ebiten.SetWindowTitle("Snake Game")
	if err := ebiten.RunGame(gameInstance); err != nil {
		log.Fatal(err)
	}
}
