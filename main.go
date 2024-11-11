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

// NewGameController crea una nueva instancia de GameController y registra el observador
func NewGameController() *GameController {
	controller := &GameController{
		loginScreen: ui.NewLoginScreen(),
	}

	controller.loginScreen.RegisterSubscriber(func(playerNames []string) {
		board := game.NewBoard(120, 60)
		snakes := []*snake.Snake{}
		for i := range playerNames {
			snakes = append(snakes, snake.NewSnake(i*2, i*2))
		}
		controller.mainGame = game.NewGame(board, snakes, playerNames)
	})

	return controller
}

// Update maneja la transición entre el login y el juego principal
func (gc *GameController) Update() error {
	// Si el juego ya se ha inicializado, actualizamos su estado
	if gc.mainGame != nil {
		return gc.mainGame.Update()
	}
	// Si el login no se ha completado, continuamos actualizando la pantalla de login
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
