package main

// Message representa un mensaje intercambiado entre cliente y servidor
type Message struct {
	Type       string   `json:"type"`       // Tipo de mensaje (e.g., "join", "move", "start_game")
	Content    string   `json:"content"`    // Contenido del mensaje (opcional, depende del tipo)
	PlayerName string   `json:"playerName"` // Nombre del jugador (opcional)
	Players    []string `json:"players"`    // Lista de jugadores en la sala de espera (opcional)
}
