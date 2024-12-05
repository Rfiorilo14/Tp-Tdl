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
	defaultWidth  = 800 // Ancho de la ventana
	defaultHeight = 600 // Alto de la ventana
	boardWidth    = 40  // Número de columnas
	boardHeight   = 30  // Número de filas
)

var (
	cellWidth  int
	cellHeight int
)

func initializeDimensions(screenWidth, screenHeight int) {
	cellWidth = screenWidth / boardWidth
	cellHeight = screenHeight / boardHeight
}

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
		initialLength := 5 // Define el tamaño inicial de la serpiente
		body := make([]Position, initialLength)
		for j := 0; j < initialLength; j++ {
			// Posiciona los segmentos consecutivos horizontalmente
			body[j] = Position{X: startingPositions[i].X - j, Y: startingPositions[i].Y}
		}

		snake := &Snake{
			Body:      body,
			Direction: "right",
			Alive:     true,
			PlayerID:  playerID,
		}
		gameState.Snakes[playerID] = snake

		// Lanzar una goroutine para manejar la serpiente
		go func(s *Snake) {
			for s.Alive {
				moveSnake(s)
				checkCollisions(s)
				time.Sleep(time.Millisecond * 200) // Ajusta la velocidad aquí
			}
		}(snake)

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

	// Colisiones con los bordes lógicos
	if head.X < 0 || head.X >= boardWidth || head.Y < 0 || head.Y >= boardHeight {
		log.Printf("Jugador %s chocó contra el borde", snake.PlayerID)
		snake.Alive = false
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
			select {
			case <-ticker.C:
				broadcastGameState() // Actualizar el estado del juego
			case msg := <-directionUpdates:
				// Procesar la dirección recibida
				updateSnakeDirection(msg)
			}
		}
	}()
}

var directionUpdates = make(chan Message, 100)

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

		// Evitar solo direcciones opuestas
		if msg.Content != oppositeDirections[snake.Direction] {
			snake.Direction = msg.Content
			log.Printf("Dirección actualizada para %s: %s", msg.PlayerName, msg.Content)
		} else {
			log.Printf("Dirección inválida para %s: %s (opuesta a %s)", msg.PlayerName, msg.Content, snake.Direction)
		}
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
