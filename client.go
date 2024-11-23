package main

import (
	"bufio"
	"fmt"
	"image/color"
	"log"
	"os"
	"strings"

	"github.com/gorilla/websocket"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font/basicfont"
)

type Game struct {
	state          string                // "waiting_room", "playing" o "scoreboard"
	waitingPlayers []string              // Jugadores en la sala de espera
	scoreboard     []string              // Tabla de puntuaciones
	conn           *websocket.Conn       // Conexión WebSocket
	snakes         map[string][]Position // Posiciones de las serpientes (por jugador)
	food           []Position            // Posiciones de la comida
}

// Update procesa la lógica del juego
func (g *Game) Update() error {
	if g.state == "waiting_room" && ebiten.IsKeyPressed(ebiten.KeyEnter) {
		// Solo envía el mensaje una vez
		err := g.conn.WriteJSON(Message{Type: "start_game"})
		if err != nil {
			log.Println("Error al enviar mensaje de inicio:", err)
		}
	} else if g.state == "scoreboard" {
		if ebiten.IsKeyPressed(ebiten.KeyR) {
			g.conn.WriteJSON(Message{Type: "restart_game"})
			g.state = "waiting_room"
		} else if ebiten.IsKeyPressed(ebiten.KeyL) {
			g.conn.WriteJSON(Message{Type: "return_to_login"})
			g.state = "waiting_room"
		}
	}
	return nil
}

// Draw dibuja los elementos en pantalla
func (g *Game) Draw(screen *ebiten.Image) {
	face := basicfont.Face7x13

	if g.state == "waiting_room" {
		text.Draw(screen, "Sala de Espera", face, 10, 20, color.White)
		y := 40
		for _, player := range g.waitingPlayers {
			text.Draw(screen, player, face, 10, y, color.White)
			y += 20
		}
		text.Draw(screen, "Presiona Enter para empezar", face, 10, y+40, color.White)

	} else if g.state == "scoreboard" {
		text.Draw(screen, "Tabla de Puntuaciones", face, 10, 20, color.White)
		y := 40
		for _, player := range g.scoreboard {
			text.Draw(screen, player, face, 10, y, color.White)
			y += 20
		}
		text.Draw(screen, "Presiona R para reiniciar", face, 10, y+40, color.White)
		text.Draw(screen, "Presiona L para volver a la sala de login", face, 10, y+60, color.White)
	} else if g.state == "playing" {
		// Dibujar el tablero
		for _, snake := range g.snakes {
			for _, segment := range snake {
				ebitenutil.DrawRect(screen, float64(segment.X*20), float64(segment.Y*20), 20, 20, color.RGBA{0, 255, 0, 255})
			}
		}

		// Dibujar comida
		for _, food := range g.food {
			ebitenutil.DrawRect(screen, float64(food.X*20), float64(food.Y*20), 20, 20, color.RGBA{255, 0, 0, 255})
		}
	}
}

// Layout define el tamaño de la pantalla del juego
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	// Definir un tamaño fijo para la ventana de juego
	return 640, 480
}

func listenToServer(conn *websocket.Conn, game *Game) {
	for {
		var msg Message
		if err := conn.ReadJSON(&msg); err != nil {
			log.Println("Error al recibir mensaje:", err)
			return
		}

		switch msg.Type {
		case "waiting_room":
			game.waitingPlayers = msg.Players
		case "start_game":
			game.state = "playing"
		case "scoreboard":
			game.scoreboard = msg.Players
			game.state = "scoreboard"
		case "game_state": // Nuevo caso para manejar el estado del juego
			if len(msg.Snakes) == 0 {
				log.Println("Error: No se recibieron serpientes en el estado del juego.")
			}
			game.snakes = make(map[string][]Position)
			for playerID, snakeBody := range msg.Snakes {
				game.snakes[playerID] = snakeBody
			}
			game.food = msg.Food

		}
	}
}

// Solicitar el nombre del jugador desde la entrada estándar
func askPlayerName() string {
	fmt.Print("Ingresa tu nombre: ")
	reader := bufio.NewReader(os.Stdin)
	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name)
	if name == "" {
		fmt.Println("El nombre no puede estar vacío. Inténtalo de nuevo.")
		return askPlayerName()
	}
	return name
}

func runClient() {
	// Pedir el nombre del jugador
	playerName := askPlayerName()

	// Conectar al servidor WebSocket
	url := "ws://localhost:8080/"
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Fatalf("Error al conectar con el servidor: %s", err)
	}
	defer conn.Close()

	// Enviar el nombre del jugador al servidor
	msg := Message{
		Type:       "join",
		PlayerName: playerName,
	}
	conn.WriteJSON(msg)

	// Inicializar el juego
	game := &Game{
		state: "waiting_room",
		conn:  conn,
	}
	go listenToServer(conn, game)

	// Ejecutar el juego
	ebiten.RunGame(game)
}

func handleInput(conn *websocket.Conn, game *Game) {
	for {
		var input string
		fmt.Scanln(&input)

		if input == "ENTER" && game.state == "waiting" {
			msg := Message{Type: "start_game"}
			if err := conn.WriteJSON(msg); err != nil {
				log.Println("Error al enviar mensaje:", err)
			}
		}
	}
}
