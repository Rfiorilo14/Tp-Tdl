package player

import pckgSnake "snake-game/snake"

type Snake = pckgSnake.Snake

type Player struct {
	name  string
	score int
	snake *Snake
	lost  bool
}

func NewPlayer(name string) *Player {
	player := &Player{

		name:  name,
		score: 0,
		snake: nil,
		lost:  false,
	}
	return player
}

func (player *Player) GetName() string {
	return player.name
}

func (player *Player) GetScore() int {
	return player.score
}

func (player *Player) Lost() bool {
	return player.lost
}

func (player *Player) SetSnake(snake *Snake) {
	player.snake = snake
}

func (player *Player) MoveSnake(direction int) {
	player.snake.Move()
	//aca se podria mover la serpiente y devolver un booleano si murio o no la serpiente y/o
	// el puntaje que obtuvo en caso de ser posible
}
