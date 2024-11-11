package ui

import (
	"strconv"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// LoginState define los estados del login
type LoginState interface {
	HandleInput(lScreen *LoginScreen)
	DrawPrompt(lScreen *LoginScreen, screen *ebiten.Image)
}

// Estado de captura de cantidad de jugadores
type AskPlayerCountState struct{}

func (s *AskPlayerCountState) HandleInput(lScreen *LoginScreen) {
	if ebiten.IsKeyPressed(ebiten.KeyEnter) {
		if num, err := parseInputToInt(lScreen.inputBuffer); err == nil && num > 0 {
			lScreen.NumPlayers = num
			lScreen.inputBuffer = ""
			lScreen.SetState(&AskPlayerNamesState{})
		}
	}
}

func (s *AskPlayerCountState) DrawPrompt(lScreen *LoginScreen, screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, "Introduce la cantidad de jugadores: "+lScreen.inputBuffer)
}

// Estado de captura de nombres de jugadores
type AskPlayerNamesState struct{}

func (s *AskPlayerNamesState) HandleInput(lScreen *LoginScreen) {
	if ebiten.IsKeyPressed(ebiten.KeyEnter) && strings.TrimSpace(lScreen.inputBuffer) != "" {
		lScreen.PlayerNames = append(lScreen.PlayerNames, lScreen.inputBuffer)
		lScreen.inputBuffer = ""

		if len(lScreen.PlayerNames) >= lScreen.NumPlayers {
			lScreen.isComplete = true
		}
	}
}

func (s *AskPlayerNamesState) DrawPrompt(lScreen *LoginScreen, screen *ebiten.Image) {
	msg := "Introduce el nombre del jugador " + strconv.Itoa(len(lScreen.PlayerNames)+1) + ": " + lScreen.inputBuffer
	ebitenutil.DebugPrint(screen, msg)
}
