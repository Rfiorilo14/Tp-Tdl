package main

import (
	"flag"
	"fmt"
	"log"
	"snake-game/client"
	"snake-game/server"
)

func main() {
	mode := flag.String("mode", "server", "Modo: 'server' para iniciar el servidor, 'client' para iniciar un cliente")
	flag.Parse()

	switch *mode {
	case "server":
		fmt.Println("Iniciando servidor...")
		server.StartServer()
	case "client":
		fmt.Println("Iniciando cliente...")
		client.StartClient()
	default:
		log.Fatalf("Modo desconocido: %s. Usa 'server' o 'client'.", *mode)
	}
}
