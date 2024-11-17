package server

import (
	"log"
	"sync"
	"time"
)

type Position struct {
	X, Y int
}

type Snake struct {
	Body      []Position
	Direction string
	PlayerID  string
	Speed     int // Velocidad de la serpiente
	Alive     bool
}

type PowerUp struct {
	Position Position
	Effect   func(*Snake)
	Duration time.Duration
}

var (
	boardWidth  = 20
	boardHeight = 20
	gameState   = struct {
		Snakes    map[string]*Snake
		Food      []Position
		Obstacles []Position
		PowerUps  []PowerUp
		mu        sync.Mutex
	}{
		Snakes:    make(map[string]*Snake),
		Food:      []Position{{X: 5, Y: 5}},
		Obstacles: []Position{{X: 10, Y: 10}},
		PowerUps:  []PowerUp{},
	}
)

func checkCollisions(snake *Snake) {
	gameState.mu.Lock()
	defer gameState.mu.Unlock()

	head := snake.Body[0]

	// Colisión con bordes
	if head.X < 0 || head.Y < 0 || head.X >= boardWidth || head.Y >= boardHeight {
		log.Printf("Jugador %s chocó contra el borde", snake.PlayerID)
		snake.Alive = false
		removePlayer(snake.PlayerID)
		return
	}

	// Colisión consigo misma
	for _, segment := range snake.Body[1:] {
		if segment == head {
			log.Printf("Jugador %s chocó consigo mismo", snake.PlayerID)
			snake.Alive = false
			removePlayer(snake.PlayerID)
			return
		}
	}

	// Colisión con otras serpientes
	for _, otherSnake := range gameState.Snakes {
		if otherSnake.PlayerID != snake.PlayerID && otherSnake.Alive {
			for _, segment := range otherSnake.Body {
				if segment == head {
					log.Printf("Jugador %s chocó con otra serpiente (%s)", snake.PlayerID, otherSnake.PlayerID)
					snake.Alive = false
					removePlayer(snake.PlayerID)
					return
				}
			}
		}
	}

	// Colisión con obstáculos
	for _, obstacle := range gameState.Obstacles {
		if obstacle == head {
			log.Printf("Jugador %s chocó contra un obstáculo", snake.PlayerID)
			snake.Alive = false
			removePlayer(snake.PlayerID)
			return
		}
	}

	// Colisión con comida
	for i, food := range gameState.Food {
		if food == head {
			log.Printf("Jugador %s comió comida", snake.PlayerID)
			tail := snake.Body[len(snake.Body)-1]
			snake.Body = append(snake.Body, tail)
			gameState.Food = append(gameState.Food[:i], gameState.Food[i+1:]...)
			return
		}
	}

	// Colisión con power-ups
	for i, powerUp := range gameState.PowerUps {
		if powerUp.Position == head {
			log.Printf("Jugador %s recogió un power-up", snake.PlayerID)
			applyPowerUp(snake, powerUp)
			gameState.PowerUps = append(gameState.PowerUps[:i], gameState.PowerUps[i+1:]...)
			break
		}
	}
}

func applyPowerUp(snake *Snake, powerUp PowerUp) {
	powerUp.Effect(snake)
	go func() {
		time.Sleep(powerUp.Duration)
		gameState.mu.Lock()
		defer gameState.mu.Unlock()
		// Revertir efecto si es necesario
		snake.Speed += 20
	}()
}

func moveSnake(snake *Snake) {
	head := snake.Body[0]
	switch snake.Direction {
	case "up":
		head.Y -= 1
	case "down":
		head.Y += 1
	case "left":
		head.X -= 1
	case "right":
		head.X += 1
	}
	snake.Body = append([]Position{head}, snake.Body[:len(snake.Body)-1]...)
}

func startGameLoop() {
	for _, snake := range gameState.Snakes {
		go func(s *Snake) {
			for s.Alive {
				moveSnake(s)
				checkCollisions(s)
				time.Sleep(time.Millisecond * time.Duration(s.Speed))
			}
		}(snake)
	}
}

func removePlayer(playerID string) {
	gameState.mu.Lock()
	defer gameState.mu.Unlock()

	if _, exists := gameState.Snakes[playerID]; exists {
		delete(gameState.Snakes, playerID)
		log.Printf("Jugador %s eliminado del juego.", playerID)
	}
}
