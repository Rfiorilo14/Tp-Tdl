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
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font/basicfont"
)

type Game struct {
	state          string   // "waiting_room", "playing" o "scoreboard"
	waitingPlayers []string // Jugadores en la sala de espera
	scoreboard     []string // Tabla de puntuaciones
	conn           *websocket.Conn
}

// Update procesa la lógica del juego
func (g *Game) Update() error {
	if g.state == "waiting_room" && ebiten.IsKeyPressed(ebiten.KeyEnter) {
		g.conn.WriteJSON(Message{Type: "start_game"})
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
