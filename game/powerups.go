// game/powerups.go
package game

import (
	"snake-game/snake"
	"time"
)

type PowerUp struct {
	SpeedBoost int
	Duration   time.Duration
}

func (p *PowerUp) ApplyToSnake(snake *snake.Snake) {
	snake.Speed += p.SpeedBoost
	go func() {
		time.Sleep(p.Duration)
		snake.Speed -= p.SpeedBoost
	}()
}
