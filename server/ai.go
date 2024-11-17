package server

import (
	"math/rand"
	"time"
)

func startAIController() {
	go func() {
		for {
			gameState.mu.Lock()
			for _, snake := range gameState.Snakes {
				if snake.PlayerID == "AI" && snake.Alive {
					directions := []string{"up", "down", "left", "right"}
					snake.Direction = directions[rand.Intn(len(directions))]
				}
			}
			gameState.mu.Unlock()
			time.Sleep(500 * time.Millisecond)
		}
	}()
}
