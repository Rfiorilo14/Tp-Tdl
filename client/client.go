package client

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"time"

	"snake-game/shared"

	"github.com/gorilla/websocket"
)

var conn *websocket.Conn
var playerName string

func StartClient() {
	// Pedir el nombre del jugador
	fmt.Print("Ingresa tu nombre: ")
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		playerName = scanner.Text()
	}
	if playerName == "" {
		log.Fatal("El nombre del jugador no puede estar vacío.")
	}

	// Conectar al servidor WebSocket
	serverAddr := "ws://127.0.0.1:8081/ws"
	var err error
	conn, _, err = websocket.DefaultDialer.Dial(serverAddr, nil)
	if err != nil {
		log.Fatalf("Error conectando al servidor WebSocket: %v", err)
		return
	}
	defer func() {
		log.Println("Cerrando conexión con el servidor...")
		conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
		conn.Close()
	}()

	log.Println("Conexión WebSocket establecida")

	// Registrar el jugador en el servidor
	err = conn.WriteJSON(shared.Message{
		Type:     "register",
		PlayerID: playerName,
	})
	if err != nil {
		log.Fatalf("Error registrando al jugador: %v", err)
		return
	}
	log.Printf("Jugador %s registrado con éxito.", playerName)

	// Iniciar el renderizado del juego
	RenderGame(conn, playerName)
}

func sendPings(conn *websocket.Conn) {
	for {
		err := conn.WriteMessage(websocket.PingMessage, nil)
		if err != nil {
			log.Printf("Error enviando ping: %v", err)
			return
		}
		time.Sleep(30 * time.Second)
	}
}

func readMessages(conn *websocket.Conn, done chan struct{}) {
	defer close(done)
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Error leyendo mensaje del servidor: %v", err)
			return
		}
		log.Printf("Mensaje del servidor: %s", message)
	}
}
