package game

import (
	"image/color"
	"log"
	"math/rand"
	"snake-game/snake"
	"snake-game/utils"
	"strconv"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const cellSize = 20

type Game struct {
	Board            *Board
	Snakes           []*snake.Snake
	CollisionManager *CollisionManager
	GameOver         bool
	Score            int
	FoodPosition     [2]int // Posición de la comida en el tablero
}

func NewGame(board *Board) *Game {
	snakes := []*snake.Snake{
		snake.NewSnake(5, 5),
	}
	game := &Game{
		Board:            board,
		Snakes:           snakes,
		CollisionManager: &CollisionManager{},
		GameOver:         false,
		Score:            0,
	}
	game.PlaceRandomFood() // Coloca la comida al iniciar el juego
	return game
}

func (g *Game) PlaceRandomFood() {
	rand.Seed(time.Now().UnixNano())
	for {
		x, y := rand.Intn(g.Board.Width), rand.Intn(g.Board.Height)
		if g.Board.Grid[y][x] == EmptyCell {
			g.FoodPosition = [2]int{x, y}
			break
		}
	}
}

func (g *Game) Update() error {
	if g.GameOver {
		return nil
	}

	for _, s := range g.Snakes {
		if ebiten.IsKeyPressed(ebiten.KeyArrowUp) && s.Direction != snake.Down {
			s.Direction = snake.Up
		} else if ebiten.IsKeyPressed(ebiten.KeyArrowDown) && s.Direction != snake.Up {
			s.Direction = snake.Down
		} else if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) && s.Direction != snake.Right {
			s.Direction = snake.Left
		} else if ebiten.IsKeyPressed(ebiten.KeyArrowRight) && s.Direction != snake.Left {
			s.Direction = snake.Right
		}
		s.Move()

		if g.CollisionManager.CheckCollisionWithBorders(g.Board, s) {
			g.GameOver = true
			log.Println("¡Perdiste! La serpiente tocó el borde del tablero.")
			return nil
		}
		if g.CollisionManager.CheckSelfCollision(s) {
			g.GameOver = true
			log.Println("¡Perdiste! La serpiente se atravesó a sí misma.")
			return nil
		}

		// Verificar colisión con la comida
		if s.Position[0] == g.FoodPosition[0] && s.Position[1] == g.FoodPosition[1] {
			s.Grow()            // Aumenta la longitud de la serpiente
			g.Score++           // Aumenta el puntaje
			g.PlaceRandomFood() // Coloca una nueva comida en el tablero
		}
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	if g.GameOver {
		msg := "¡Perdiste! Puntaje: " + strconv.Itoa(g.Score)
		ebitenutil.DebugPrintAt(screen, msg, (g.Board.Width*cellSize)/2-80, (g.Board.Height*cellSize)/2)
		return
	}

	screen.Fill(color.Black)

	for y := range g.Board.Grid {
		for x := range g.Board.Grid[y] {
			g.Board.Grid[y][x] = EmptyCell
		}
	}

	for _, s := range g.Snakes {
		for _, pos := range s.Body {
			if g.Board.isInBounds(pos[0], pos[1]) {
				g.Board.Grid[pos[1]][pos[0]] = SnakeCell
			}
		}
	}

	g.Board.Grid[g.FoodPosition[1]][g.FoodPosition[0]] = FoodCell

	for y, row := range g.Board.Grid {
		for x, cell := range row {
			var cellColor color.Color
			switch cell {
			case SnakeCell:
				cellColor = utils.ColSnake
			case FoodCell:
				cellColor = utils.ColFood
			case ObstacleCell:
				cellColor = utils.ColObstacle
			default:
				continue
			}
			ebitenutil.DrawRect(screen, float64(x*cellSize), float64(y*cellSize), cellSize, cellSize, cellColor)
		}
	}

	// Simular un puntaje más grande dibujando el texto varias veces con desplazamiento
	scoreText := "Puntaje: " + strconv.Itoa(g.Score)
	xPos, yPos := 20, 40

	// Crear un efecto de texto grande superponiéndolo con pequeños desplazamientos
	for dx := -1; dx <= 1; dx++ {
		for dy := -1; dy <= 1; dy++ {
			ebitenutil.DebugPrintAt(screen, scoreText, xPos+dx, yPos+dy)
		}
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return g.Board.Width * cellSize, g.Board.Height * cellSize
}
