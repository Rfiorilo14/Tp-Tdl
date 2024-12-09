# Snake Multiplayer Game  
**Autor:** Roy Fiorilo  
**Legajo:** 108419  

## 📖 Descripción del Proyecto  
Este proyecto implementa un juego multijugador inspirado en el clásico *Snake*, desarrollado en **Go**. Los jugadores pueden conectarse a un servidor, competir en tiempo real, y visualizar los resultados en una tabla de puntuaciones. La arquitectura incluye un servidor centralizado que gestiona el estado del juego y clientes que interactúan mediante WebSockets.  

## ✨ Funcionalidades  
- 🕹️ **Sala de espera** para jugadores antes de iniciar una partida.  
- 🚀 **Juego en tiempo real** con movimiento de serpientes y generación dinámica de comida.  
- ❌ **Detección de colisiones** (paredes, cuerpo propio y otras serpientes).  
- 🏆 **Tabla de puntuaciones** al final de cada partida.  
- 🔁 **Reinicio de partidas** y retorno a la sala de espera.  

## 📂 Estructura del Proyecto  
El proyecto está dividido en los siguientes archivos:  

### 1. `client.go`  
Contiene la lógica del cliente, incluida la interfaz gráfica utilizando **Ebiten**. Permite a los jugadores conectarse al servidor, interactuar con el juego y enviar comandos como movimientos.  

### 2. `game_logic.go`  
Define la lógica central del juego. Aquí se implementa el movimiento de las serpientes, la detección de colisiones y la generación de comida. También se manejan las estructuras principales como `Snake` y `Position`.  

### 3. `server.go`  
Gestión del servidor utilizando **Gorilla WebSocket**. Maneja conexiones, sincroniza el estado global del juego y envía actualizaciones a los clientes.  

### 4. `types.go`  
Define las estructuras de datos compartidas entre cliente y servidor, como `Message`, que facilita la comunicación a través de JSON.  

### 5. `main.go`  
Archivo de entrada del programa. Permite ejecutar el proyecto en modo servidor o cliente, según los argumentos proporcionados.  

## 🛠️ Tecnologías Utilizadas  
- **Lenguaje:** Go  
- **Bibliotecas:**  
  - [**Gorilla WebSocket**](https://github.com/gorilla/websocket): Para la comunicación en tiempo real.  
  - [**Ebiten**](https://ebiten.org/): Para la interfaz gráfica del cliente.  
- **JSON**: Para estructurar y transmitir mensajes entre cliente y servidor.  
