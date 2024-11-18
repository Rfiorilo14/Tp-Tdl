package client

import (
	"fmt"
	"image/color"
	"log"
	"snake-game/shared"

	"github.com/gorilla/websocket"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type SnakeSegment struct {
	X, Y int
}

type Game struct {
	State        shared.GameState
	PlayerName   string
	Conn         *websocket.Conn
	stateChannel chan shared.GameState // Canal para actualizar el estado del juego
	playerSnake  []SnakeSegment        // Serpiente del jugador
	direction    string                // Dirección actual de movimiento
	boardWidth   int                   // Ancho del tablero
	boardHeight  int                   // Altura del tablero
	tick         int                   // Control de velocidad de movimiento
}

func (g *Game) startListening() {
	go func() {
		for {
			var newState shared.GameState
			err := g.Conn.ReadJSON(&newState)
			if err != nil {
				log.Printf("Error leyendo del servidor: %v", err)
				close(g.stateChannel) // Cerrar el canal al desconectarse
				return
			}
			g.stateChannel <- newState // Enviar el nuevo estado al canal
		}
	}()
}

func (g *Game) Update() error {
	// Actualización periódica para el movimiento continuo
	g.tick++
	if g.tick%10 == 0 { // Ajusta este valor para controlar la velocidad
		g.moveSnake()
	}

	// Leer el estado desde el canal, no directamente del servidor
	select {
	case newState, ok := <-g.stateChannel:
		if !ok {
			return fmt.Errorf("conexión cerrada por el servidor")
		}
		g.State = newState
	default:
		// No hay nuevos estados, continuar normalmente
	}

	// Manejar entrada del jugador (actualizar dirección)
	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) && g.direction != "down" {
		g.direction = "up"
	} else if ebiten.IsKeyPressed(ebiten.KeyArrowDown) && g.direction != "up" {
		g.direction = "down"
	} else if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) && g.direction != "right" {
		g.direction = "left"
	} else if ebiten.IsKeyPressed(ebiten.KeyArrowRight) && g.direction != "left" {
		g.direction = "right"
	}

	return nil
}

func (g *Game) moveSnake() {
	// Determinar el nuevo segmento de la cabeza de la serpiente
	head := g.playerSnake[0]
	var newHead SnakeSegment

	switch g.direction {
	case "up":
		newHead = SnakeSegment{X: head.X, Y: head.Y - 1}
	case "down":
		newHead = SnakeSegment{X: head.X, Y: head.Y + 1}
	case "left":
		newHead = SnakeSegment{X: head.X - 1, Y: head.Y}
	case "right":
		newHead = SnakeSegment{X: head.X + 1, Y: head.Y}
	}

	// Asegurarse de que la serpiente no salga del tablero
	if newHead.X < 0 {
		newHead.X = g.boardWidth - 1
	} else if newHead.X >= g.boardWidth {
		newHead.X = 0
	}
	if newHead.Y < 0 {
		newHead.Y = g.boardHeight - 1
	} else if newHead.Y >= g.boardHeight {
		newHead.Y = 0
	}

	// Insertar la nueva cabeza y eliminar el último segmento
	g.playerSnake = append([]SnakeSegment{newHead}, g.playerSnake[:len(g.playerSnake)-1]...)
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Dibujar el tablero y los elementos del juego

	// Dibujar la serpiente del jugador
	for _, segment := range g.playerSnake {
		ebitenutil.DrawRect(screen, float64(segment.X*20), float64(segment.Y*20), 20, 20, color.RGBA{R: 0, G: 255, B: 0, A: 255})
	}

	// Dibujar otras serpientes (si es necesario)
	for playerID, snake := range g.State.Snakes {
		if playerID == g.PlayerName {
			continue // No dibujar nuestra propia serpiente
		}
		clr := color.RGBA{R: 255, G: 255, B: 0, A: 255} // Serpientes de otros jugadores
		for _, segment := range snake.Body {
			ebitenutil.DrawRect(screen, float64(segment.X*20), float64(segment.Y*20), 20, 20, clr)
		}
	}

	// Dibujar la comida
	for _, food := range g.State.Food {
		ebitenutil.DrawRect(screen, float64(food.X*20), float64(food.Y*20), 20, 20, color.RGBA{R: 255, G: 0, B: 0, A: 255})
	}

	// Mostrar mensajes en pantalla
	ebitenutil.DebugPrint(screen, fmt.Sprintf("Jugador: %s", g.PlayerName))
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return 1200, 600 // Dimensiones del tablero
}

func (g *Game) sendMove(direction string) {
	// Enviar el movimiento al servidor
	err := g.Conn.WriteJSON(shared.Message{
		Type:      "move",
		PlayerID:  g.PlayerName,
		Direction: direction,
	})
	if err != nil {
		log.Printf("Error enviando movimiento: %v", err)
	}
}

func RenderGame(conn *websocket.Conn, playerName string) {
	// Configurar el estado inicial del juego
	game := &Game{
		State:        shared.GameState{},
		PlayerName:   playerName,
		Conn:         conn,
		stateChannel: make(chan shared.GameState, 1),
		playerSnake:  []SnakeSegment{{X: 10, Y: 10}}, // Inicializar la serpiente en el centro
		direction:    "right",                        // Dirección inicial
		boardWidth:   1200 / 20,                      // Calcular ancho en segmentos
		boardHeight:  600 / 20,                       // Calcular altura en segmentos
	}

	// Expandir la serpiente inicial
	for i := 1; i < 10; i++ {
		game.playerSnake = append(game.playerSnake, SnakeSegment{X: 10 - i, Y: 10})
	}

	// Iniciar la escucha de mensajes del servidor
	game.startListening()

	// Configurar ventana del juego
	ebiten.SetWindowSize(1200, 600)
	ebiten.SetWindowTitle("Juego de la Viborita - Cliente")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
