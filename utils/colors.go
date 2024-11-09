// utils/colors.go
package utils

import "image/color"

var (
	ColSnake    = color.NRGBA{R: 0, G: 255, B: 0, A: 255} // Verde para la serpiente
	ColFood     = color.NRGBA{R: 255, G: 0, B: 0, A: 255} // Rojo para la comida
	ColObstacle = color.NRGBA{R: 0, G: 0, B: 255, A: 255} // Azul para los obstáculos
)
