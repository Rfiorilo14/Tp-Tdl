// snake/snake.go
package snake

import "image/color"

type Snake struct {
	id            int
	bodyPositions [][2]int // Guarda las posiciones del cuerpo
	direction     string   // Dirección de movimiento actual
	speed         int
	alive         bool
	color         *color.RGBA
}

var directions = map[string][2]int{
	"up":    {0, -1},
	"down":  {0, 1},
	"left":  {-1, 0},
	"right": {1, 0},
}

/*
var (
	Up    = [2]int{0, -1}
	Down  = [2]int{0, 1}
	Left  = [2]int{-1, 0}
	Right = [2]int{1, 0}
)

El usuario ingresa en el cliente una tecla. Esa tecla pasa a ser un token. Ese token puede
ser "up","down","left","right".
*/

// NewSnake crea una nueva instancia de la serpiente en una posición dada
func NewSnake(x, y int, id int) *Snake {
	snake := &Snake{
		id:            id,
		bodyPositions: [][2]int{{x, y}},
		direction:     "",
		speed:         1,
		alive:         true,
	}
	return snake
}

func (snake *Snake) SetInitialDirection(direction string) {
	snake.direction = direction
}

func (snake *Snake) SetColor(color *color.RGBA) {
	snake.color = color
}

func (snake *Snake) GetColor() *color.RGBA {
	return snake.color
}

func (snake *Snake) GetDirection() string {
	return snake.direction
}

/*
Podemos hacer que en caso de que el jugador quiera una direccion que ya esta establecida en su
snake, se ejecute esta funcion pero no va a cambiar nada
*/
func (snake *Snake) Move(direction string) {

}

/*
// Grow aumenta la longitud de la serpiente al comer
func (s *Snake) Grow() {
	s.Length++
}
*/
