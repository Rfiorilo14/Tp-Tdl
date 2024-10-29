Alumnos:
Roy Fiorilo 
Brian Fernandez 
Maximiliano Torre

Idea de la Estructura del Proyecto :

|-- main.go                   // Punto de entrada del juego
|-- game                      // Lógica principal del juego
    |-- game.go               // Control de flujo del juego, lógica general
    |-- board.go              // Lógica del tablero, generación de comida y obstáculos
    |-- collision.go          // Detección de colisiones
    |-- powerups.go           // Implementación de power-ups y power-downs
|-- snake                     // Lógica de las serpientes
    |-- snake.go              // Clase/struct de serpiente, movimientos, velocidad
    |-- ia.go                 // IA para las serpientes controladas por la computadora
|-- ui                        // Interfaz gráfica y renderizado
    |-- render.go             // Renderizado 3D, carga de modelos y texturas
    |-- blender_import.go     // Funciones para importar y manejar recursos desde Blender
|-- assets                    // Recursos de Blender y gráficos
    |-- models                // Modelos 3D exportados de Blender (ej: .obj, .fbx)
    |-- textures              // Texturas para los modelos 3D
    |-- animations            // Archivos de animación, si los usas
|-- models                    // Modelos de datos para jugadores, tablero, etc.
|-- utils                     // Utilidades como generación aleatoria, temporizadores
