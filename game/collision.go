// game/collision.go
package game

import (
	"snake-game/snake"
)

type CollisionManager struct{}

// CheckCollisionWithBorders verifica si la serpiente ha colisionado con los bordes del tablero
func (c *CollisionManager) CheckCollisionWithBorders(board *Board, s *snake.Snake) bool {
	// La posición de la cabeza ahora se encuentra en s.Body[0]
	x, y := s.Body[0][0], s.Body[0][1]
	return x < 0 || x >= board.Width || y < 0 || y >= board.Height
}

// CheckSelfCollision verifica si la serpiente se toca a sí misma
func (c *CollisionManager) CheckSelfCollision(s *snake.Snake) bool {
	for i := 1; i < len(s.Body); i++ { // Comienza en 1 para evitar la cabeza misma
		if s.Body[i] == s.Body[0] { // La cabeza colisiona con el cuerpo
			return true
		}
	}
	return false
}
