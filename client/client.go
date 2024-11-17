package client

import (
	"bufio"
	"log"
	"os"
	"time"

	"github.com/gorilla/websocket"
)

var conn *websocket.Conn

func StartClient() {
	serverAddr := "ws://127.0.0.1:8081/ws"
	conn, _, err := websocket.DefaultDialer.Dial(serverAddr, nil)
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

	done := make(chan struct{})

	// Inicia pings para mantener la conexión activa
	go sendPings(conn)

	// Inicia recepción de mensajes del servidor
	go readMessages(conn, done)

	// Canal para manejar la entrada del usuario
	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			text := scanner.Text()
			if text == "exit" {
				log.Println("Cerrando conexión por solicitud del usuario...")
				conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "Cliente cerrando conexión"))
				conn.Close()
				close(done)
				os.Exit(0)
			}

			// Enviar mensaje al servidor
			err := conn.WriteMessage(websocket.TextMessage, []byte(text))
			if err != nil {
				log.Printf("Error enviando mensaje: %v", err)
				break
			}
		}
	}()

	// Esperar hasta que se cierre el canal done
	<-done
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

func reconnect() {
	log.Println("Reconectando en 5 segundos...")
	if conn != nil {
		conn.Close() // Cierra la conexión previa si existe
	}
	time.Sleep(5 * time.Second) // Espera antes de reconectar
	StartClient()               // Reinicia el cliente
}

func sendInputs(conn *websocket.Conn) {
	// Aquí envías inputs o mensajes al servidor
	err := conn.WriteJSON(map[string]string{"action": "move", "direction": "up"})
	if err != nil {
		log.Printf("Error enviando input: %v", err)
	}
}
