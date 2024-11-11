// main.go
package main

import (
	"log"
	"snake-game/game"
	"snake-game/snake"
	"snake-game/ui"

	"github.com/hajimehoshi/ebiten/v2"
)

type GameController struct {
	loginScreen *ui.LoginScreen
	mainGame    *game.Game
}

func NewGameController() *GameController {
	return &GameController{
		loginScreen: ui.NewLoginScreen(),
	}
}

// Update maneja la transición entre el login y el juego principal
func (gc *GameController) Update() error {
	// Verifica si el login está completo y, de ser así, inicializa el juego principal
	if gc.mainGame == nil {
		playerNames, isComplete := gc.loginScreen.StartGame()
		if isComplete {
			board := game.NewBoard(120, 60)
			snakes := []*snake.Snake{}
			for i := range playerNames {
				log.Printf("Jugador %d: %s", i+1, playerNames[i])
				snakes = append(snakes, snake.NewSnake(i*2, i*2)) // Crea serpiente para cada jugador
			}
			gc.mainGame = game.NewGame(board)
			gc.mainGame.Snakes = snakes
		}
	}

	// Si el juego está inicializado, llama a su Update, si no, llama al login
	if gc.mainGame != nil {
		return gc.mainGame.Update()
	}
	return gc.loginScreen.Update()
}

// Draw maneja el dibujo de la pantalla actual (login o juego principal)
func (gc *GameController) Draw(screen *ebiten.Image) {
	if gc.mainGame != nil {
		gc.mainGame.Draw(screen)
	} else {
		gc.loginScreen.Draw(screen)
	}
}

// Layout define las dimensiones de la ventana
func (gc *GameController) Layout(outsideWidth, outsideHeight int) (int, int) {
	if gc.mainGame != nil {
		return gc.mainGame.Layout(outsideWidth, outsideHeight)
	}
	return gc.loginScreen.Layout(outsideWidth, outsideHeight)
}

func main() {
	// Configuración de ventana
	ebiten.SetWindowSize(1200, 600)
	ebiten.SetWindowTitle("Juego de la Viborita Multijugador")

	// Inicializa el controlador principal del juego
	gameController := NewGameController()

	// Ejecuta el controlador principal que maneja tanto el login como el juego
	if err := ebiten.RunGame(gameController); err != nil {
		log.Fatal(err)
	}
}
