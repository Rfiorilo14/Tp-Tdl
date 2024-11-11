package ui

import "github.com/hajimehoshi/ebiten/v2"

// Screen define la interfaz de todas las pantallas
type Screen interface {
	Update() error
	Draw(screen *ebiten.Image)
	Layout(outsideWidth, outsideHeight int) (int, int)
}

// ScreenFactory crea pantallas específicas
type ScreenFactory struct{}

func (f *ScreenFactory) CreateScreen(screenType string) Screen {
	switch screenType {
	case "login":
		return NewLoginScreen()
	// Aquí se podrían agregar más pantallas (como la pantalla de juego)
	default:
		return nil
	}
}
