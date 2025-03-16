package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"server/internal/server"
	"server/internal/server/clients"
)

var (
	port = flag.Int("port", 8080, "port to listen on")
)

func main() {
	flag.Parse()

	//define game hub
	hub := server.NewHub()

	//define handler for websocket connections
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		hub.Serve(clients.NewWebSocketClient, w, r)
	})

	go hub.Run()
	addr := fmt.Sprintf(":%d", *port)
	log.Printf("Starting Server on %s", addr)
	err := http.ListenAndServe(addr, nil)

	if err != nil {
		log.Fatal("Failed to start server: ", err)
	}

}
