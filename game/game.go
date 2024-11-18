package game

import (
	"fmt"
	"image/color"
	"math/rand"
	"snake-game/snake"
	"snake-game/utils"
	"sort"
	"strconv"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const cellSize = 20

type Game struct {
	Board       *Board
	Snakes      []*snake.Snake
	PlayerNames []string // Almacena nombres de los jugadores

	CollisionManager *CollisionManager //mmm nose
	GameOver         bool
	Score            int
}

func NewGame(board *Board, snakes []*snake.Snake, playerNames []string) *Game {

	game := &Game{
		Board:       board,
		Snakes:      snakes,
		PlayerNames: playerNames,

		CollisionManager: &CollisionManager{},
		GameOver:         false,
		Score:            0,
	}
	game.PlaceRandomFood()
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

/*
func (g *Game) Update() error {
	if g.GameOver {
		return nil
	}

	activeSnakes := 0

	for i, s := range g.Snakes {
		if !s.Alive {
			continue // Ignorar serpientes que ya no están vivas
		}

		// Comprobar si la serpiente ha comido la comida
		if s.Position[0] == g.FoodPosition[0] && s.Position[1] == g.FoodPosition[1] {
			s.Grow()
			g.Score++
			s.Score++
			g.PlaceRandomFood()
		}

		// Usa la estrategia de control correspondiente para cada serpiente
		g.ControlStrategies[i].UpdateDirection(s)
		s.Move()

		// Verificar colisiones
		if g.CollisionManager.CheckCollisionWithBorders(g.Board, s) {
			log.Printf("¡La serpiente %d tocó el borde y ha sido eliminada!", i+1)
			s.Alive = false
		}
		if g.CollisionManager.CheckSelfCollision(s) {
			log.Printf("¡La serpiente %d se atravesó a sí misma y ha sido eliminada!", i+1)
			s.Alive = false
		}

		// Contar cuántas serpientes siguen activas
		if s.Alive {
			activeSnakes++
		}
	}

	// Si no quedan serpientes activas, terminar el juego
	if activeSnakes == 0 {
		g.GameOver = true
		log.Println("¡Todas las serpientes han sido eliminadas! Fin del juego.")
	}

	return nil
}
*/

// TODO ESTO PARECE QUE PERTENECE A LA INTERFAZ
func (g *Game) Draw(screen *ebiten.Image) {
	if g.GameOver {
		g.ShowRankedTable(screen)
		return
	}

	screen.Fill(color.Black)
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

	// Dibujar solo las serpientes vivas
	for _, s := range g.Snakes {
		if !s.Alive {
			continue
		}
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

	scoreText := "Puntaje: " + strconv.Itoa(g.Score)
	xPos, yPos := 20, 40

	for dx := -1; dx <= 1; dx++ {
		for dy := -1; dy <= 1; dy++ {
			ebitenutil.DebugPrintAt(screen, scoreText, xPos+dx, yPos+dy)
		}
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return g.Board.Width * cellSize, g.Board.Height * cellSize
}

func (g *Game) ShowRankedTable(screen *ebiten.Image) {
	type ScoreEntry struct {
		Name  string
		Score int
	}

	// Crear una lista de puntajes
	var scoreEntries []ScoreEntry
	for i, s := range g.Snakes {
		scoreEntries = append(scoreEntries, ScoreEntry{Name: g.PlayerNames[i], Score: s.Score})
	}

	// Ordenar por puntaje en orden descendente
	sort.Slice(scoreEntries, func(i, j int) bool {
		return scoreEntries[i].Score > scoreEntries[j].Score
	})

	// Construir el mensaje de la tabla de puntajes
	msg := "Ranking Final:\n"
	for i, entry := range scoreEntries {
		msg += fmt.Sprintf("%d. %s - Puntaje: %d\n", i+1, entry.Name, entry.Score)
	}

	// Posiciones iniciales para el texto en pantalla
	xPos, yPos := 20, 20

	// Dibujar el texto con un "borde" para simular mayor tamaño
	for dx := -1; dx <= 1; dx++ {
		for dy := -1; dy <= 1; dy++ {
			if dx != 0 || dy != 0 { // Evitar el centro, que será el texto principal
				ebitenutil.DebugPrintAt(screen, msg, xPos+dx, yPos+dy)
			}
		}
	}

	// Dibujar el texto principal en blanco encima para mejorar la visibilidad
	ebitenutil.DebugPrintAt(screen, msg, xPos, yPos)
}
