package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	context := newContext()

	go context.run()

	http.Handle("/", http.FileServer(http.Dir("./client/build")))

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(context, w, r)
	})

	port := ":" + os.Getenv("PORT")
	log.Fatal(http.ListenAndServe(port, nil))
}
