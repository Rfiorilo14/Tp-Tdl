package main

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

// ServerState representa el estado global del servidor.
type ServerState struct {
	players        map[string]*Player
	eliminated     []*Player
	gameStarted    bool
	currentPlayers int
	mu             sync.Mutex
}

var serverState = &ServerState{
	players:        make(map[string]*Player),
	eliminated:     []*Player{},
	gameStarted:    false,
	currentPlayers: 0,
}

// Player representa a un jugador conectado.
type Player struct {
	Name string
	Conn *websocket.Conn
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// handleConnections maneja nuevas conexiones WebSocket.
func handleConnections(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error al actualizar conexión:", err)
		return
	}
	defer conn.Close()

	var playerName string

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
			startGameIfPossible()

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

		case "update_direction":
			updateSnakeDirection(msg)
		}
	}
}

// startGameIfPossible inicia el juego si las condiciones son correctas.
func startGameIfPossible() {
	serverState.mu.Lock()
	defer serverState.mu.Unlock()

	if serverState.currentPlayers > 0 && !serverState.gameStarted {
		log.Println("Iniciando el juego...")
		serverState.gameStarted = true
		initializeGame()
		startGameLoop()

		for _, player := range serverState.players {
			err := player.Conn.WriteJSON(Message{Type: "start_game"})
			if err != nil {
				log.Printf("Error al enviar mensaje de inicio a %s: %s", player.Name, err)
			}
		}
	} else {
		log.Printf("No se puede iniciar el juego. Jugadores: %d Juego iniciado: %t", serverState.currentPlayers, serverState.gameStarted)
	}
}

// broadcastWaitingRoom envía la lista de jugadores en la sala de espera a todos los clientes.
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

// startGame notifica a los jugadores que la partida ha comenzado.
func startGame() {
	for _, player := range serverState.players {
		err := player.Conn.WriteJSON(Message{
			Type: "start_game",
		})
		if err != nil {
			log.Printf("Error al enviar mensaje de inicio a %s: %s", player.Name, err)
		}
	}
	log.Println("El juego ha comenzado!")
}

// checkEndGame verifica si todos los jugadores están eliminados.
func checkEndGame() {
	serverState.mu.Lock()
	defer serverState.mu.Unlock()

	if serverState.currentPlayers == 0 {
		broadcastScoreboard()
		serverState.gameStarted = false
	}
}

// broadcastScoreboard envía la tabla de puntuaciones a todos los clientes.
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

// restartGame reinicia la partida con los mismos jugadores.
func restartGame() {
	serverState.mu.Lock()
	defer serverState.mu.Unlock()

	serverState.eliminated = []*Player{}
	serverState.gameStarted = true

	initializeGame()
	startGameLoop()
}

// resetToLogin regresa a todos los jugadores a la sala de login.
func resetToLogin() {
	serverState.mu.Lock()
	defer serverState.mu.Unlock()

	serverState.eliminated = []*Player{}
	serverState.gameStarted = false

	broadcastWaitingRoom()
}
