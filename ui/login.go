// Aplicamos el patrón Observer y Singleton
package ui

import (
	"strconv"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
)

type LoginScreen struct {
	NumPlayers  int
	PlayerNames []string
	inputBuffer string
	isComplete  bool
	state       LoginState       // Estado actual del login
	subscribers []func([]string) // Observadores
}

var instance *LoginScreen

// NewLoginScreen crea una instancia única de LoginScreen
func NewLoginScreen() *LoginScreen {
	if instance == nil {
		instance = &LoginScreen{
			PlayerNames: []string{},
			state:       &AskPlayerCountState{},
		}
	}
	return instance
}

func (l *LoginScreen) SetState(state LoginState) {
	l.state = state
}

func (l *LoginScreen) RegisterSubscriber(subscriber func([]string)) {
	l.subscribers = append(l.subscribers, subscriber)
}

func (l *LoginScreen) notifySubscribers() {
	for _, subscriber := range l.subscribers {
		subscriber(l.PlayerNames)
	}
}

func (l *LoginScreen) Update() error {
	if l.isComplete {
		l.notifySubscribers()
		return nil
	}

	for _, r := range ebiten.InputChars() {
		if r != '\n' && r != '\r' {
			l.inputBuffer += string(r)
		}
	}

	l.state.HandleInput(l)

	if ebiten.IsKeyPressed(ebiten.KeyBackspace) && len(l.inputBuffer) > 0 {
		l.inputBuffer = l.inputBuffer[:len(l.inputBuffer)-1]
	}

	return nil
}

func (l *LoginScreen) Draw(screen *ebiten.Image) {
	l.state.DrawPrompt(l, screen)
}

func (l *LoginScreen) Layout(outsideWidth, outsideHeight int) (int, int) {
	return 800, 600
}

func parseInputToInt(input string) (int, error) {
	input = strings.TrimSpace(input)
	return strconv.Atoi(input)
}
