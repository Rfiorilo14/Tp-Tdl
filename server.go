package main

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

// Estado del servidor
type ServerState struct {
	players        map[string]*Player // Jugadores conectados
	eliminated     []*Player          // Jugadores eliminados
	gameStarted    bool               // Indica si el juego comenzó
	currentPlayers int                // Número de jugadores en la partida
	mu             sync.Mutex         // Mutex para manejar concurrencia
}

// Inicializa el estado del servidor
var serverState = &ServerState{
	players:        make(map[string]*Player),
	eliminated:     []*Player{},
	gameStarted:    false,
	currentPlayers: 0,
}

// Estructura del jugador
type Player struct {
	Name string
	Conn *websocket.Conn
}

// Manejador de WebSocket
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Permitir todas las conexiones
	},
}

// Maneja nuevas conexiones WebSocket
func handleConnections(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error al actualizar conexión:", err)
		return
	}

	defer conn.Close()

	var playerName string

	// Procesar mensajes del cliente
	for {
		var msg Message
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Println("Error al leer mensaje:", err)
			break
		}

		switch msg.Type {
		case "join":
			serverState.mu.Lock()
			if _, exists := serverState.players[msg.PlayerName]; !exists {
				playerName = msg.PlayerName
				serverState.players[playerName] = &Player{
					Name: playerName,
					Conn: conn,
				}
				serverState.currentPlayers++
			}
			serverState.mu.Unlock()
			broadcastWaitingRoom()

		case "start_game":
			if !serverState.gameStarted {
				serverState.gameStarted = true
				startGame()
			}

		case "player_eliminated":
			serverState.mu.Lock()
			if player, exists := serverState.players[msg.PlayerName]; exists {
				serverState.eliminated = append(serverState.eliminated, player)
				delete(serverState.players, msg.PlayerName)
				serverState.currentPlayers--
			}
			serverState.mu.Unlock()
			checkEndGame()

		case "restart_game":
			restartGame()

		case "return_to_login":
			resetToLogin()
		}
	}
}

// Difundir la sala de espera a todos los jugadores
func broadcastWaitingRoom() {
	serverState.mu.Lock()
	defer serverState.mu.Unlock()

	players := []string{}
	for name := range serverState.players {
		players = append(players, name)
	}

	for _, player := range serverState.players {
		err := player.Conn.WriteJSON(Message{
			Type:    "waiting_room",
			Players: players,
		})
		if err != nil {
			log.Printf("Error al enviar mensaje a %s: %s", player.Name, err)
		}
	}
}

// Iniciar la partida
func startGame() {
	for _, player := range serverState.players {
		player.Conn.WriteJSON(Message{
			Type: "start_game",
		})
	}
}

// Verifica si todos los jugadores están eliminados
func checkEndGame() {
	if serverState.currentPlayers == 0 {
		broadcastScoreboard()
		serverState.gameStarted = false
	}
}

// Difundir la tabla de puntuaciones
func broadcastScoreboard() {
	scoreboard := []string{}
	for _, player := range serverState.eliminated {
		scoreboard = append(scoreboard, player.Name)
	}

	for _, player := range serverState.eliminated {
		err := player.Conn.WriteJSON(Message{
			Type:    "scoreboard",
			Players: scoreboard,
		})
		if err != nil {
			log.Printf("Error al enviar mensaje a %s: %s", player.Name, err)
		}
	}
}

// Reiniciar la partida con los mismos jugadores
func restartGame() {
	serverState.mu.Lock()
	serverState.eliminated = []*Player{}
	serverState.gameStarted = true
	serverState.mu.Unlock()

	startGame()
}

// Regresar a la sala de login
func resetToLogin() {
	serverState.mu.Lock()
	serverState.eliminated = []*Player{}
	serverState.gameStarted = false
	serverState.mu.Unlock()

	broadcastWaitingRoom()
}