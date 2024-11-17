package client

import (
	"log"
	"time"

	"github.com/gorilla/websocket"
)

var conn *websocket.Conn

func StartClient() {
	var err error
	conn, _, err = websocket.DefaultDialer.Dial("ws://localhost:8080/ws", nil)
	if err != nil {
		log.Fatal("Error conectando al servidor:", err)
	}
	defer conn.Close()

	done := make(chan struct{})
	go readMessages(done)
	go sendPings() // Nueva función para mantener la conexión activa

	sendInputs()
}

func sendPings() {
	for {
		err := conn.WriteMessage(websocket.PingMessage, []byte{})
		if err != nil {
			log.Println("Error enviando ping al servidor:", err)
			return
		}
		time.Sleep(30 * time.Second) // Intervalo para enviar pings
	}
}

func readMessages(done chan struct{}) {
	defer close(done)
	for {
		var msg map[string]interface{}
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Println("Conexión perdida. Intentando reconectar...")
			reconnect()
			break
		}
		log.Println("Mensaje del servidor:", msg)
	}
}

func reconnect() {
	time.Sleep(5 * time.Second) // Espera antes de reconectar
	StartClient()
}

func sendInputs() {
	// Ejemplo de mensaje enviado al servidor
	msg := map[string]string{
		"type":      "move",
		"player_id": "player1",
		"direction": "up",
	}
	conn.WriteJSON(msg)
}
