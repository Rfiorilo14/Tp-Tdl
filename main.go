package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

func main() {
	// Define el modo de operación: servidor o cliente
	mode := flag.String("mode", "server", "Modo de operación: server o client")
	flag.Parse()

	if *mode == "server" {
		fmt.Println("Iniciando servidor...")
		http.HandleFunc("/", handleConnections)

		port := ":8080"
		fmt.Println("Servidor escuchando en", port)
		log.Fatal(http.ListenAndServe(port, nil))
	} else if *mode == "client" {
		runClient()
	} else {
		log.Fatalf("Modo desconocido: %s", *mode)
	}
}
