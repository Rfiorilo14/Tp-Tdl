package server

import (
	"context"
	"log"
	"net"
	"net/http"
	"sync"
	"time"

	"snake-game/shared"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Permitir todas las conexiones, para pruebas
	},
}

var (
	clients   = make(map[*websocket.Conn]string)
	broadcast = make(chan shared.Message)
	mu        sync.Mutex
	server    *http.Server // Servidor HTTP para permitir cierre controlado
)

func HandleConnection(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Error actualizando a WebSocket: %v", err)
		return
	}
	defer func() {
		log.Println("Cerrando conexión con el cliente:", r.RemoteAddr)
		conn.Close()
	}()

	log.Println("Jugador conectado:", r.RemoteAddr)

	for {
		// Leer mensaje del cliente
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Error inesperado: %v", err)
			}
			break
		}

		log.Printf("Mensaje recibido: %s", message)

		// Responder al cliente
		response := "Mensaje recibido correctamente: " + string(message)
		err = conn.WriteMessage(messageType, []byte(response))
		if err != nil {
			log.Printf("Error al enviar mensaje: %v", err)
			break
		}
	}
}

func registerPlayer(conn *websocket.Conn, playerID string) {
	mu.Lock()
	defer mu.Unlock()
	clients[conn] = playerID
}

func handleMessages(conn *websocket.Conn, playerID string) {
	defer func() {
		mu.Lock()
		delete(clients, conn)
		mu.Unlock()
		log.Printf("Jugador %s desconectado (ID Conexión: %v).", playerID, conn.RemoteAddr())
		conn.Close()
	}()

	for {
		var msg shared.Message
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Printf("Error al leer mensaje de %s (ID Conexión: %v): %v", playerID, conn.RemoteAddr(), err)
			break
		}
		broadcast <- msg
	}
}

func handleBroadcasts() {
	for {
		msg := <-broadcast
		mu.Lock()
		for conn, playerID := range clients {
			if err := conn.WriteJSON(msg); err != nil {
				log.Printf("Error enviando mensaje a %s: %v. Eliminando conexión.", playerID, err)
				conn.Close()
				delete(clients, conn)
			}
		}
		mu.Unlock()
	}
}

func StartServer() {
	// Cambia el puerto aquí
	server := &http.Server{Addr: "127.0.0.1:8081"}

	listener, err := net.Listen("tcp", server.Addr)
	if err != nil {
		log.Fatalf("Error iniciando el servidor: %v", err)
	}

	// Ajustar configuración de Keep-Alive para el listener
	tcpListener := listener.(*net.TCPListener)
	listener = tcpKeepAliveListener{tcpListener}

	http.HandleFunc("/ws", HandleConnection)
	go handleBroadcasts()

	log.Println("Servidor iniciado en :8081") // Asegúrate de actualizar el log con el nuevo puerto
	if err := server.Serve(listener); err != http.ErrServerClosed {
		log.Fatalf("Error en el servidor: %v", err)
	}
}

// tcpKeepAliveListener configura opciones de Keep-Alive para conexiones entrantes
type tcpKeepAliveListener struct {
	*net.TCPListener
}

func (ln tcpKeepAliveListener) Accept() (net.Conn, error) {
	tc, err := ln.TCPListener.AcceptTCP()
	if err != nil {
		return nil, err
	}
	tc.SetKeepAlive(true)
	tc.SetKeepAlivePeriod(20 * time.Second) // Ajusta el intervalo de Keep-Alive
	return tc, nil
}

func shutdownServer(ctx context.Context) {
	timeoutCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Detener el servidor HTTP
	if err := server.Shutdown(timeoutCtx); err != nil {
		log.Fatalf("Error cerrando el servidor: %v", err)
	}

	// Cerrar todas las conexiones WebSocket activas
	mu.Lock()
	for conn := range clients {
		conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "Servidor apagado"))
		conn.Close()
		delete(clients, conn)
	}
	mu.Unlock()

	log.Println("Servidor cerrado de forma controlada.")
}
