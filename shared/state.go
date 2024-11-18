package shared

import (
	"image/color"
)

type Position struct {
	X, Y int
}

type Snake struct {
	Body      []Position
	Direction string
	PlayerID  string
	Alive     bool
	Color     color.RGBA // Cambio a color.RGBA
}

type PowerUp struct {
	Position Position
	Effect   func(*Snake)
	Duration int
}

type GameState struct {
	Snakes    map[string]*Snake
	Food      []Position
	Obstacles []Position
	PowerUps  []PowerUp
}
