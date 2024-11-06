Alumnos:
Roy Fiorilo 
Brian Fernandez 
Maximiliano Torre


## Idea de la Estructura del Proyecto

```plaintext
|-- go.mod                   // Archivo de módulos de Go, define dependencias y versión del proyecto
|-- main.go                  // Punto de partida del juego, aca ejecutamos
|-- game                     // Lógica principal del juego
|   |-- game.go              // Control de flujo del juego, inicialización y lógica general
|   |-- board.go             // Lógica del tablero, generación de comida y obstáculos
|   |-- collision.go         // Detección de colisiones entre la serpiente y otros objetos
|   |-- powerups.go          // Implementación de power-ups y power-downs
|-- snake                    // Lógica de las serpientes (jugador y posibles oponentes)
|   |-- snake.go             // Estructura de la serpiente, control de movimientos y velocidad
|   |-- ia.go                // IA para serpientes controladas por la computadora
|-- ui                       // Interfaz gráfica y renderizado del juego
|   |-- render.go            // Renderizado 3D, carga de modelos y texturas
|   |-- blender_import.go    // Funciones para importar y manejar recursos desde Blender
|-- assets                   // Recursos gráficos y modelos 3D
|   |-- models               // Modelos 3D exportados de Blender (ej: .obj, .fbx)
|   |-- textures             // Texturas para los modelos 3D
|   |-- animations           // Archivos de animación para el juego
|-- models                   // Modelos de datos para estructurar información como jugadores, tablero, etc.
|-- utils                    // Utilidades varias, como generación aleatoria, temporizadores
|   |-- colors               // Colores para la serpiente, comida, etc
