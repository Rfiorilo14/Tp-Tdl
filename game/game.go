// game/game.go
package game

import (
	"image/color"
	"log"
	"snake-game/snake"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const cellSize = 20

type Game struct {
	Board            *Board
	Snakes           []*snake.Snake
	CollisionManager *CollisionManager
	GameOver         bool
	Score            int // Puntaje acumulado
}

// NewGame crea una nueva instancia del juego
func NewGame(board *Board) *Game {
	snakes := []*snake.Snake{
		snake.NewSnake(5, 5, "Player1"),
	}
	return &Game{
		Board:            board,
		Snakes:           snakes,
		CollisionManager: &CollisionManager{},
		GameOver:         false,
		Score:            0, // Puntaje inicializado en cero
	}
}

// Update actualiza la lógica del juego
func (g *Game) Update() error {
	if g.GameOver {
		return nil // Detenemos el juego si está en estado de GameOver
	}

	for _, s := range g.Snakes {
		// Captura de entrada del teclado para cambiar la dirección
		if ebiten.IsKeyPressed(ebiten.KeyUp) {
			s.ChangeDirection(snake.Up)
		} else if ebiten.IsKeyPressed(ebiten.KeyDown) {
			s.ChangeDirection(snake.Down)
		} else if ebiten.IsKeyPressed(ebiten.KeyLeft) {
			s.ChangeDirection(snake.Left)
		} else if ebiten.IsKeyPressed(ebiten.KeyRight) {
			s.ChangeDirection(snake.Right)
		}

		// Mover la serpiente
		s.Move()

		// Verificar colisiones con los bordes
		if g.CollisionManager.CheckCollisionWithBorders(g.Board, s) {
			g.GameOver = true
			log.Println("¡Perdiste! La serpiente tocó el borde del tablero.")
			return nil
		}

		// Verificar colisión consigo misma
		if g.CollisionManager.CheckSelfCollision(s) {
			g.GameOver = true
			log.Println("¡Perdiste! La serpiente se tocó a sí misma.")
			return nil
		}

		// Verificar colisión con comida
		if g.Board.Grid[s.Position[1]][s.Position[0]] == FoodCell {
			s.Grow()                  // La serpiente crece al comer
			g.Score++                 // Aumenta el puntaje
			g.Board.PlaceRandomFood() // Coloca una nueva comida en el tablero
		}

		// Actualizar la posición en el tablero
		for _, pos := range s.Body {
			g.Board.PlaceSnake(pos[0], pos[1])
		}
	}
	return nil
}

// Draw dibuja el tablero, la serpiente y el puntaje en la pantalla
func (g *Game) Draw(screen *ebiten.Image) {
	if g.GameOver {
		msg := "¡Perdiste! Puntaje: " + strconv.Itoa(g.Score)
		ebitenutil.DebugPrintAt(screen, msg, (g.Board.Width*cellSize)/2-40, (g.Board.Height*cellSize)/2)
		return
	}

	// Dibuja el tablero y la serpiente en la pantalla
	for y, row := range g.Board.Grid {
		for x, cell := range row {
			var cellColor color.Color
			switch cell {
			case SnakeCell:
				cellColor = color.NRGBA{R: 0, G: 255, B: 0, A: 255} // Verde para la serpiente
			case FoodCell:
				cellColor = color.NRGBA{R: 255, G: 0, B: 0, A: 255} // Rojo para la comida
			case ObstacleCell:
				cellColor = color.NRGBA{R: 0, G: 0, B: 255, A: 255} // Azul para obstáculos
			default:
				cellColor = color.NRGBA{R: 0, G: 0, B: 0, A: 255} // Negro para celdas vacías
			}
			// Dibuja el rectángulo de la celda
			ebitenutil.DrawRect(screen, float64(x*cellSize), float64(y*cellSize), cellSize, cellSize, cellColor)
		}
	}

	// Mostrar puntaje en la esquina superior izquierda
	ebitenutil.DebugPrintAt(screen, "Puntaje: "+strconv.Itoa(g.Score), 10, 10)
}

// Layout especifica el tamaño de la pantalla
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return g.Board.Width * cellSize, g.Board.Height * cellSize
}
