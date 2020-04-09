package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
	"github.com/peterzernia/set/game"
)

// Keeps track of all connections
var clients = make(map[*websocket.Conn]bool)

// Context contains necessary clients throughout the application
type Context struct {
	Game *game.Game
}

func main() {
	context := Context{}
	hub := newHub()

	go hub.run(&context)

	http.Handle("/", http.FileServer(http.Dir("./client/build")))

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(&context, hub, w, r)
	})

	port := ":" + os.Getenv("PORT")
	fmt.Println("Running server on port " + port)
	log.Fatal(http.ListenAndServe(port, nil))
}
