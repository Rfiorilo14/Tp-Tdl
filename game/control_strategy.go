package game

import (
	"snake-game/snake"

	"github.com/hajimehoshi/ebiten/v2"
)

type ControlStrategy interface {
	UpdateDirection(s *snake.Snake)
}

type ArrowControlStrategy struct{}

func (a *ArrowControlStrategy) UpdateDirection(s *snake.Snake) {
	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) && s.Direction != snake.Down {
		s.Direction = snake.Up
	} else if ebiten.IsKeyPressed(ebiten.KeyArrowDown) && s.Direction != snake.Up {
		s.Direction = snake.Down
	} else if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) && s.Direction != snake.Right {
		s.Direction = snake.Left
	} else if ebiten.IsKeyPressed(ebiten.KeyArrowRight) && s.Direction != snake.Left {
		s.Direction = snake.Right
	}
}

type WASDControlStrategy struct{}

func (w *WASDControlStrategy) UpdateDirection(s *snake.Snake) {
	if ebiten.IsKeyPressed(ebiten.KeyW) && s.Direction != snake.Down {
		s.Direction = snake.Up
	} else if ebiten.IsKeyPressed(ebiten.KeyS) && s.Direction != snake.Up {
		s.Direction = snake.Down
	} else if ebiten.IsKeyPressed(ebiten.KeyA) && s.Direction != snake.Right {
		s.Direction = snake.Left
	} else if ebiten.IsKeyPressed(ebiten.KeyD) && s.Direction != snake.Left {
		s.Direction = snake.Right
	}
}

type KOLÑControlStrategy struct{}

func (k *KOLÑControlStrategy) UpdateDirection(s *snake.Snake) {
	if ebiten.IsKeyPressed(ebiten.KeyO) && s.Direction != snake.Down {
		s.Direction = snake.Up
	} else if ebiten.IsKeyPressed(ebiten.KeyL) && s.Direction != snake.Up {
		s.Direction = snake.Down
	} else if ebiten.IsKeyPressed(ebiten.KeyK) && s.Direction != snake.Right {
		s.Direction = snake.Left
	} else if ebiten.IsKeyPressed(ebiten.KeySemicolon) && s.Direction != snake.Left {
		s.Direction = snake.Right
	}
}
