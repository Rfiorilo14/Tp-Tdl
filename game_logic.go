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
	Body      []Position // Cuerpo de la serpiente
	Direction string     // Dirección actual ("up", "down", "left", "right")
	PlayerID  string     // ID del jugador propietario
	Alive     bool       // Si la serpiente está activa
}

var (
	boardWidth  = 20
	boardHeight = 20
	gameState   = struct {
		Snakes map[string]*Snake
		Food   []Position
		mu     sync.Mutex
	}{
		Snakes: make(map[string]*Snake),
		Food:   []Position{{X: 5, Y: 5}},
	}
)

func initializeGame() {
	gameState.mu.Lock()
	defer gameState.mu.Unlock()

	// Crear serpientes para los jugadores
	startingPositions := []Position{
		{X: 3, Y: 3},
		{X: boardWidth - 4, Y: 3},
		{X: 3, Y: boardHeight - 4},
	}
	i := 0
	for playerID := range serverState.players {
		gameState.Snakes[playerID] = &Snake{
			Body:      []Position{startingPositions[i]},
			Direction: "right",
			Alive:     true,
			PlayerID:  playerID,
		}
		i++
		if i >= len(startingPositions) {
			startingPositions = append(startingPositions, Position{X: rand.Intn(boardWidth), Y: rand.Intn(boardHeight)})
		}
	}

	// Garantizar que haya comida inicial
	if len(gameState.Food) == 0 {
		spawnFood()
	}
	log.Printf("Estado inicial del juego: %+v", gameState)
}

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

	// Insertar la nueva cabeza
	snake.Body = append([]Position{newHead}, snake.Body...)

	// Eliminar la cola si no creció
	snake.Body = snake.Body[:len(snake.Body)-1]
}

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

	for i, food := range gameState.Food {
		if food == head {
			log.Printf("Jugador %s comió comida", snake.PlayerID)

			// Agregar un nuevo segmento a la serpiente
			tail := snake.Body[len(snake.Body)-1]
			snake.Body = append(snake.Body, tail)

			// Eliminar comida y generar una nueva
			gameState.Food = append(gameState.Food[:i], gameState.Food[i+1:]...)
			spawnFood()
			return
		}
	}
}
func spawnFood() {
	foodPosition := Position{
		X: rand.Intn(boardWidth),
		Y: rand.Intn(boardHeight),
	}

	// Asegurarse de que no haya colisión con una serpiente u obstáculo
	for _, snake := range gameState.Snakes {
		for _, segment := range snake.Body {
			if segment == foodPosition {
				spawnFood() // Intentar de nuevo
				return
			}
		}
	}

	gameState.Food = append(gameState.Food, foodPosition)
}

var gameBoard = struct {
	Grid   [][]CellType      // Representación de la matriz del tablero
	Snakes map[string]*Snake // Serpientes asociadas a jugadores
	Food   Position          // Posición de la comida
	mu     sync.Mutex        // Mutex para concurrencia
}{
	Grid:   make([][]CellType, boardHeight),
	Snakes: make(map[string]*Snake),
}

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

func startGameLoop() {
	initializeGame()
	go func() {
		ticker := time.NewTicker(200 * time.Millisecond)
		defer ticker.Stop()

		broadcastGameState() // Difundir estado inicial

		for serverState.gameStarted {
			<-ticker.C
			updateGameState()    // Actualizar posiciones y colisiones
			broadcastGameState() // Difundir el estado actualizado
		}
	}()
}

func removePlayer(playerID string) {
	gameState.mu.Lock()
	defer gameState.mu.Unlock()

	if _, exists := gameState.Snakes[playerID]; exists {
		delete(gameState.Snakes, playerID)
		log.Printf("Jugador %s eliminado del juego.", playerID)
	}
}

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
