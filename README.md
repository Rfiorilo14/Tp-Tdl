Alumnos:
Roy Fiorilo 
Brian Fernandez 
Maximiliano Torre


## Estructura del Proyecto

```plaintext
|-- go.mod                   // Archivo de módulos de Go, define dependencias y versión del proyecto
|-- main.go                  // Punto de entrada del juego, donde se ejecuta la inicialización
|-- game                     // Lógica principal del juego
|   |-- board.go             // Lógica del tablero, generación de comida y obstáculos
|   |-- collision.go         // Detección de colisiones entre la serpiente y otros objetos
|   |-- control_strategy.go  // Estrategias de control para manejar las acciones de las serpientes
|   |-- game.go              // Control del flujo del juego, inicialización y lógica general
|   |-- powerups.go          // Implementación de power-ups y power-downs en el juego
|-- snake                    // Lógica específica de las serpientes (jugador y posibles oponentes)
|   |-- snake.go             // Estructura y lógica de la serpiente, control de movimientos y velocidad
|-- ui                       // Interfaz gráfica y manejo de pantallas del juego
|   |-- login_state.go       // Estado y lógica de la pantalla de login
|   |-- login.go             // Lógica de autenticación o introducción de usuarios
|   |-- screen_factory.go    // Creación y gestión de diferentes pantallas de UI del juego
|-- utils                    // Utilidades varias para el juego
|   |-- colors.go            // Definición de colores utilizados en la interfaz del juego y serpientes
