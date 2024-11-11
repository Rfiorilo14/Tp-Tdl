// ui/login.go
package ui

import (
	"strconv"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// LoginScreen representa la pantalla de inicio de sesión para capturar la cantidad de jugadores y sus nombres.
type LoginScreen struct {
	NumPlayers  int
	PlayerNames []string
	inputBuffer string
	currentStep int // 0: pedir cantidad de jugadores, 1: pedir nombres
	isComplete  bool
}

// NewLoginScreen crea una nueva instancia de LoginScreen.
func NewLoginScreen() *LoginScreen {
	return &LoginScreen{
		NumPlayers:  0,
		PlayerNames: []string{},
		inputBuffer: "",
		currentStep: 0,
		isComplete:  false,
	}
}

// Update maneja la lógica de actualización de la pantalla de inicio de sesión.
func (l *LoginScreen) Update() error {
	if l.isComplete {
		return nil // No hace nada si el login ya se ha completado.
	}

	// Captura las entradas de texto.
	for _, r := range ebiten.InputChars() {
		if r != '\n' && r != '\r' {
			l.inputBuffer += string(r) // Agrega caracteres al buffer de entrada.
		}
	}

	// Captura el evento de presionar la tecla Enter.
	if ebiten.IsKeyPressed(ebiten.KeyEnter) {
		if l.currentStep == 0 {
			// Paso para capturar la cantidad de jugadores.
			if num, err := parseInputToInt(l.inputBuffer); err == nil && num > 0 {
				l.NumPlayers = num
				l.currentStep = 1
				l.inputBuffer = "" // Limpia el buffer para el siguiente paso.
			}
		} else if l.currentStep == 1 && len(l.PlayerNames) < l.NumPlayers {
			// Paso para capturar los nombres de los jugadores.
			if strings.TrimSpace(l.inputBuffer) != "" {
				l.PlayerNames = append(l.PlayerNames, l.inputBuffer)
				l.inputBuffer = "" // Limpia el buffer para el siguiente nombre.

				// Verifica si ya se capturaron todos los nombres.
				if len(l.PlayerNames) >= l.NumPlayers {
					l.isComplete = true
				}
			}
		}
	}

	// Captura la tecla Backspace para borrar caracteres.
	if ebiten.IsKeyPressed(ebiten.KeyBackspace) && len(l.inputBuffer) > 0 {
		l.inputBuffer = l.inputBuffer[:len(l.inputBuffer)-1]
	}

	return nil
}

// Draw dibuja la pantalla de inicio de sesión.
func (l *LoginScreen) Draw(screen *ebiten.Image) {
	if l.currentStep == 0 {
		ebitenutil.DebugPrint(screen, "Introduce la cantidad de jugadores: "+l.inputBuffer)
	} else if l.currentStep == 1 && len(l.PlayerNames) < l.NumPlayers {
		msg := "Introduce el nombre del jugador " + strconv.Itoa(len(l.PlayerNames)+1) + ": " + l.inputBuffer
		ebitenutil.DebugPrint(screen, msg)
	}
}

// Layout establece el tamaño de la pantalla de inicio de sesión.
func (l *LoginScreen) Layout(outsideWidth, outsideHeight int) (int, int) {
	return 800, 600 // Dimensiones de la ventana.
}

// StartGame devuelve los nombres de los jugadores y verifica si la pantalla de inicio ha terminado.
func (l *LoginScreen) StartGame() ([]string, bool) {
	if l.isComplete {
		return l.PlayerNames, true
	}
	return nil, false
}

// Función auxiliar para convertir la entrada de texto en un número entero.
func parseInputToInt(input string) (int, error) {
	input = strings.TrimSpace(input)
	return strconv.Atoi(input)
}
