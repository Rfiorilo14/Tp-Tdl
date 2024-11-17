package client

import (
	"image/color"
	"log"
	"snake-game/shared"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Game struct {
	State shared.GameState
}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	for _, snake := range g.State.Snakes {
		for _, segment := range snake.Body {
			// Dibuja cada segmento con un color específico
			ebitenutil.DrawRect(screen, float64(segment.X*20), float64(segment.Y*20), 20, 20, color.RGBA{R: 0, G: 255, B: 0, A: 255})
		}
	}
	ebitenutil.DebugPrint(screen, "Juego de la Viborita en Go")
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return 640, 480
}

func RenderGame() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Juego de la Viborita - Cliente")
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
