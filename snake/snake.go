// snake/snake.go
package snake

var (
	Up    = [2]int{0, -1}
	Down  = [2]int{0, 1}
	Left  = [2]int{-1, 0}
	Right = [2]int{1, 0}
)

type Snake struct {
	ID        int
	Position  [2]int
	Body      [][2]int // Guarda todas las posiciones del cuerpo de la serpiente
	Direction [2]int
	Speed     int
	Alive     bool
	Length    int // Longitud actual de la serpiente
}

// NewSnake crea una nueva instancia de la serpiente en una posición dada
func NewSnake(x, y int, name string) *Snake {
	return &Snake{
		Position:  [2]int{x, y},
		Body:      [][2]int{{x, y}}, // Inicializamos el cuerpo con la posición inicial
		Direction: Right,
		Speed:     1,
		Alive:     true,
		Length:    5, // Longitud inicial de la serpiente
	}
}

// ChangeDirection cambia la dirección de la serpiente si no es la dirección opuesta
func (s *Snake) ChangeDirection(newDirection [2]int) {
	if s.Direction[0] != -newDirection[0] || s.Direction[1] != -newDirection[1] {
		s.Direction = newDirection
	}
}

// Move actualiza la posición de la serpiente y su cuerpo
func (s *Snake) Move() {
	// Calcula la nueva posición de la cabeza
	newHead := [2]int{s.Position[0] + s.Direction[0], s.Position[1] + s.Direction[1]}

	// Inserta la nueva posición al inicio del cuerpo
	s.Body = append([][2]int{newHead}, s.Body...)

	// Si la longitud del cuerpo excede la longitud actual de la serpiente, eliminamos el último segmento
	if len(s.Body) > s.Length {
		s.Body = s.Body[:len(s.Body)-1]
	}

	// Actualizamos la posición de la cabeza
	s.Position = newHead
}

// IsSelfCollision verifica si la cabeza de la serpiente colisiona con su propio cuerpo
func (s *Snake) IsSelfCollision() bool {
	// Comparamos la posición de la cabeza con el resto del cuerpo
	for i := 1; i < len(s.Body); i++ {
		if s.Body[i] == s.Position {
			return true
		}
	}
	return false
}

// Grow aumenta la longitud de la serpiente
func (s *Snake) Grow() {
	s.Length++
}
