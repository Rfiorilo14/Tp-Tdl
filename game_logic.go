// game_logic.go
package main

import (
	"log"
	"math/rand"
	"sync"
	"time"
)

type CellType int

const (
	Empty CellType = iota
	Food
)

type Position struct {
	X, Y int
}

type Snake struct {
	Body      []Position
	Direction string
	PlayerID  string
	Alive     bool
}

var (
	gameState = struct {
		Snakes map[string]*Snake
		Food   []Position
		mu     sync.Mutex
	}{
		Snakes: make(map[string]*Snake),
		Food:   []Position{{X: 5, Y: 5}},
	}
)

const (
	defaultWidth  = 800
	defaultHeight = 600
	boardWidth    = 40
	boardHeight   = 30
)

var (
	cellWidth  int
	cellHeight int
)

// initializeDimensions calcula las dimensiones de las celdas basado en el tamaño del tablero.
func initializeDimensions(screenWidth, screenHeight int) {
	cellWidth = screenWidth / boardWidth
	cellHeight = screenHeight / boardHeight
}

// initializeGame inicializa el estado del juego, creando serpientes y comida inicial.
func initializeGame() {
	gameState.mu.Lock()
	defer gameState.mu.Unlock()

	startingPositions := []Position{
		{X: 3, Y: 3},
		{X: boardWidth - 4, Y: 3},
		{X: 3, Y: boardHeight - 4},
	}
	i := 0
	for playerID := range serverState.players {
		initialLength := 5
		body := make([]Position, initialLength)
		for j := 0; j < initialLength; j++ {
			body[j] = Position{X: startingPositions[i].X - j, Y: startingPositions[i].Y}
		}

		snake := &Snake{
			Body:      body,
			Direction: "right",
			Alive:     true,
			PlayerID:  playerID,
		}
		gameState.Snakes[playerID] = snake

		go func(s *Snake) {
			for s.Alive {
				moveSnake(s)
				checkCollisions(s)
				time.Sleep(time.Millisecond * 200)
			}
		}(snake)

		i++
		if i >= len(startingPositions) {
			startingPositions = append(startingPositions, Position{X: rand.Intn(boardWidth), Y: rand.Intn(boardHeight)})
		}
	}

	if len(gameState.Food) == 0 {
		spawnFood()
	}
	log.Printf("Estado inicial del juego: %+v", gameState)
}

// moveSnake mueve la serpiente en la dirección actual y ajusta su posición.
func moveSnake(snake *Snake) {
	if !snake.Alive {
		return
	}

	head := snake.Body[0]
	var newHead Position

	switch snake.Direction {
	case "up":
		newHead = Position{X: head.X, Y: head.Y - 1}
	case "down":
		newHead = Position{X: head.X, Y: head.Y + 1}
	case "left":
		newHead = Position{X: head.X - 1, Y: head.Y}
	case "right":
		newHead = Position{X: head.X + 1, Y: head.Y}
	}

	snake.Body = append([]Position{newHead}, snake.Body...)
	snake.Body = snake.Body[:len(snake.Body)-1]
}

// checkCollisions verifica si una serpiente ha colisionado y maneja las consecuencias.
func checkCollisions(snake *Snake) {
	gameState.mu.Lock()
	defer gameState.mu.Unlock()

	head := snake.Body[0]

	if head.X < 0 || head.X >= boardWidth || head.Y < 0 || head.Y >= boardHeight {
		log.Printf("Jugador %s chocó contra el borde", snake.PlayerID)
		snake.Alive = false
		return
	}

	for _, segment := range snake.Body[1:] {
		if segment == head {
			log.Printf("Jugador %s chocó consigo mismo", snake.PlayerID)
			snake.Alive = false
			removePlayer(snake.PlayerID)
			return
		}
	}

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

	for i, food := range gameState.Food {
		if food == head {
			log.Printf("Jugador %s comió comida", snake.PlayerID)
			tail := snake.Body[len(snake.Body)-1]
			snake.Body = append(snake.Body, tail)

			gameState.Food = append(gameState.Food[:i], gameState.Food[i+1:]...)
			spawnFood()
			return
		}
	}
}

// spawnFood genera comida en una posición aleatoria del tablero.
func spawnFood() {
	foodPosition := Position{
		X: rand.Intn(boardWidth),
		Y: rand.Intn(boardHeight),
	}

	for _, snake := range gameState.Snakes {
		for _, segment := range snake.Body {
			if segment == foodPosition {
				spawnFood()
				return
			}
		}
	}

	gameState.Food = append(gameState.Food, foodPosition)
}

// updateGameState actualiza el estado global del juego.
func updateGameState() {
	gameState.mu.Lock()
	defer gameState.mu.Unlock()

	for _, snake := range gameState.Snakes {
		if snake.Alive {
			moveSnake(snake)
			checkCollisions(snake)
		}
	}
}

// startGameLoop inicia el bucle principal del juego.
func startGameLoop() {
	initializeGame()

	go func() {
		ticker := time.NewTicker(200 * time.Millisecond)
		defer ticker.Stop()

		broadcastGameState()

		for serverState.gameStarted {
			select {
			case <-ticker.C:
				broadcastGameState()
			case msg := <-directionUpdates:
				updateSnakeDirection(msg)
			}
		}
	}()
}

var directionUpdates = make(chan Message, 100)

// updateSnakeDirection actualiza la dirección de una serpiente según los mensajes del cliente.
func updateSnakeDirection(msg Message) {
	gameState.mu.Lock()
	defer gameState.mu.Unlock()

	if snake, exists := gameState.Snakes[msg.PlayerName]; exists && snake.Alive {
		oppositeDirections := map[string]string{
			"up":    "down",
			"down":  "up",
			"left":  "right",
			"right": "left",
		}

		if msg.Content != oppositeDirections[snake.Direction] {
			log.Printf("Actualizando dirección para %s: %s -> %s", msg.PlayerName, snake.Direction, msg.Content)
			snake.Direction = msg.Content
		} else {
			log.Printf("Dirección inválida para %s: %s (opuesta a %s)", msg.PlayerName, msg.Content, snake.Direction)
		}
	}
}

// removePlayer elimina un jugador del juego.
func removePlayer(playerID string) {
	gameState.mu.Lock()
	defer gameState.mu.Unlock()

	if _, exists := gameState.Snakes[playerID]; exists {
		delete(gameState.Snakes, playerID)
		log.Printf("Jugador %s eliminado del juego.", playerID)
	}
}

// broadcastGameState difunde el estado actual del juego a todos los jugadores conectados.
func broadcastGameState() {
	gameState.mu.Lock()
	defer gameState.mu.Unlock()

	state := Message{
		Type:   "game_state",
		Snakes: make(map[string][]Position),
		Food:   gameState.Food,
	}

	for playerID, snake := range gameState.Snakes {
		if snake.Alive {
			state.Snakes[playerID] = snake.Body
		}
	}

	log.Printf("Estado a enviar: %+v", state)

	for playerID, player := range serverState.players {
		if player.Conn == nil {
			removePlayer(playerID)
			continue
		}
		err := player.Conn.WriteJSON(state)
		if err != nil {
			log.Printf("Error enviando estado a %s: %v", playerID, err)
			removePlayer(playerID)
		}
	}
}
