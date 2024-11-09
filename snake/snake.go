// snake/snake.go
package snake

type Snake struct {
	ID        int
	Position  [2]int   // Posición de la cabeza
	Body      [][2]int // Guarda las posiciones del cuerpo
	Direction [2]int   // Dirección de movimiento actual
	Speed     int
	Alive     bool
	Length    int // Longitud actual de la serpiente
}

var (
	Up    = [2]int{0, -1}
	Down  = [2]int{0, 1}
	Left  = [2]int{-1, 0}
	Right = [2]int{1, 0}
)

// NewSnake crea una nueva instancia de la serpiente en una posición dada
func NewSnake(x, y int) *Snake {
	return &Snake{
		Position:  [2]int{x, y},
		Body:      [][2]int{{x, y}}, // Inicializa el cuerpo en la posición inicial
		Direction: Right,            // Ahora "Right" estará correctamente definida
		Speed:     1,
		Alive:     true,
		Length:    6, // Longitud inicial
	}
}

func (s *Snake) Move() {
	// Calcula la nueva posición de la cabeza
	newHead := [2]int{s.Position[0] + s.Direction[0], s.Position[1] + s.Direction[1]}

	// Inserta la nueva posición al inicio del cuerpo
	s.Body = append([][2]int{newHead}, s.Body...)

	// Mantiene la longitud si no ha comido (cuando el tamaño excede la longitud permitida)
	if len(s.Body) > s.Length {
		s.Body = s.Body[:len(s.Body)-1] // Elimina el último segmento
	}

	// Actualizamos la posición de la cabeza
	s.Position = newHead
}

// Grow aumenta la longitud de la serpiente al comer
func (s *Snake) Grow() {
	s.Length++
}
