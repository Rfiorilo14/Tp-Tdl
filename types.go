package main

// Message representa un mensaje intercambiado entre cliente y servidor
type Message struct {
	Type       string                `json:"type"`       // Tipo de mensaje
	Content    string                `json:"content"`    // Contenido adicional
	PlayerName string                `json:"playerName"` // Nombre del jugador
	Players    []string              `json:"players"`    // Lista de jugadores
	Snakes     map[string][]Position `json:"snakes"`     // Estado de las serpientes
	Food       []Position            `json:"food"`       // Posiciones de la comida
}
