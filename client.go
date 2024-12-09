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
	state          string
	waitingPlayers []string
	scoreboard     []string
	conn           *websocket.Conn
	playerName     string
	snakes         map[string][]Position
	food           []Position
	lastDirection  string
}

// Update procesa la lógica del juego, incluyendo teclas y cambios de estado.
func (g *Game) Update() error {
	if g.state == "waiting_room" && ebiten.IsKeyPressed(ebiten.KeyEnter) {
		err := g.conn.WriteJSON(Message{Type: "start_game"})
		if err != nil {
			log.Println("Error al enviar mensaje de inicio:", err)
		}
	} else if g.state == "playing" {
		var newDirection string
		if ebiten.IsKeyPressed(ebiten.KeyArrowUp) || ebiten.IsKeyPressed(ebiten.KeyW) {
			newDirection = "up"
		} else if ebiten.IsKeyPressed(ebiten.KeyArrowDown) || ebiten.IsKeyPressed(ebiten.KeyS) {
			newDirection = "down"
		} else if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) || ebiten.IsKeyPressed(ebiten.KeyA) {
			newDirection = "left"
		} else if ebiten.IsKeyPressed(ebiten.KeyArrowRight) || ebiten.IsKeyPressed(ebiten.KeyD) {
			newDirection = "right"
		}

		if newDirection != "" && newDirection != g.lastDirection {
			log.Printf("Intentando enviar dirección: %s (última: %s)", newDirection, g.lastDirection)
			err := g.conn.WriteJSON(Message{
				Type:       "update_direction",
				PlayerName: g.playerName,
				Content:    newDirection,
			})
			if err != nil {
				log.Println("Error al enviar nueva dirección:", err)
			} else {
				g.lastDirection = newDirection
			}
		}
	}
	return nil
}

// Draw dibuja los elementos visuales del juego según el estado actual.
func (g *Game) Draw(screen *ebiten.Image) {
	face := basicfont.Face7x13
	initializeDimensions(640, 480)
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
		for _, snake := range g.snakes {
			for _, segment := range snake {
				x := float64(segment.X * cellWidth)
				y := float64(segment.Y * cellHeight)
				ebitenutil.DrawRect(screen, x, y, float64(cellWidth), float64(cellHeight), color.RGBA{0, 255, 0, 255})
			}
		}

		for _, food := range g.food {
			x := float64(food.X * cellWidth)
			y := float64(food.Y * cellHeight)
			ebitenutil.DrawRect(screen, x, y, float64(cellWidth), float64(cellHeight), color.RGBA{255, 0, 0, 255})
		}
	}
}

// Layout ajusta el tamaño de la pantalla según las dimensiones externas.
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	initializeDimensions(outsideWidth, outsideHeight)
	return outsideWidth, outsideHeight
}

// listenToServer escucha mensajes del servidor y actualiza el estado del juego.
func listenToServer(conn *websocket.Conn, game *Game) {
	for {
		var msg Message
		if err := conn.ReadJSON(&msg); err != nil {
			log.Println("Error al recibir mensaje:", err)
			continue
		}
		switch msg.Type {
		case "waiting_room":
			game.waitingPlayers = msg.Players
		case "start_game":
			game.state = "playing"
		case "game_state":
			game.snakes = make(map[string][]Position)
			for playerID, snakeBody := range msg.Snakes {
				game.snakes[playerID] = snakeBody
			}
			game.food = msg.Food
		case "scoreboard":
			game.state = "scoreboard"
			game.scoreboard = msg.Players
		}
	}
}

// askPlayerName solicita al usuario que ingrese su nombre para usar en el juego.
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

// runClient inicializa el cliente, conecta al servidor y ejecuta el juego.
func runClient() {
	playerName := askPlayerName()
	url := "ws://localhost:8080/"
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Fatalf("Error al conectar con el servidor: %s", err)
	}
	defer conn.Close()

	msg := Message{
		Type:       "join",
		PlayerName: playerName,
	}
	conn.WriteJSON(msg)

	game := &Game{
		state:      "waiting_room",
		conn:       conn,
		playerName: playerName,
	}
	go listenToServer(conn, game)
	ebiten.RunGame(game)
}

// handleInput procesa entradas del usuario para interactuar con el servidor.
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

// drawGrid dibuja la cuadrícula del tablero en la pantalla.
func drawGrid(screen *ebiten.Image) {
	for x := 0; x < boardWidth; x++ {
		for y := 0; y < boardHeight; y++ {
			rectX := float64(x * cellWidth)
			rectY := float64(y * cellHeight)
			ebitenutil.DrawRect(screen, rectX, rectY, float64(cellWidth), float64(cellHeight), color.RGBA{0, 0, 0, 255})
		}
	}
}
