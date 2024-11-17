package server

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"snake-game/shared"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

var (
	clients   = make(map[*websocket.Conn]string)
	broadcast = make(chan shared.Message)
	mu        sync.Mutex
)

func HandleConnection(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error al conectar:", err)
		return
	}
	defer conn.Close()

	playerID := fmt.Sprintf("player%d", len(clients)+1)
	registerPlayer(conn, playerID)

	log.Println("Jugador conectado:", playerID)
	go handleMessages(conn, playerID)
}

func registerPlayer(conn *websocket.Conn, playerID string) {
	mu.Lock()
	defer mu.Unlock()
	clients[conn] = playerID
}

func handleMessages(conn *websocket.Conn, playerID string) {
	defer func() {
		mu.Lock()
		delete(clients, conn)
		mu.Unlock()
		log.Printf("Jugador %s desconectado.", playerID)
	}()

	for {
		var msg shared.Message
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Printf("Error al leer mensaje de %s: %v", playerID, err)
			break
		}
		broadcast <- msg
	}
}

func handleBroadcasts() {
	for {
		msg := <-broadcast
		mu.Lock()
		for conn := range clients {
			if err := conn.WriteJSON(msg); err != nil {
				log.Println("Error enviando mensaje:", err)
				conn.Close()
				delete(clients, conn)
			}
		}
		mu.Unlock()
	}
}

func StartServer() {
	http.HandleFunc("/ws", HandleConnection)
	go handleBroadcasts()
	log.Fatal(http.ListenAndServe(":8080", nil))
}
