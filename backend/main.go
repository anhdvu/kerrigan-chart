package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"

	"github.com/gorilla/websocket"
)

var clients = make(map[*websocket.Conn]bool)
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func home(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Hey yo!")
	w.Write([]byte("Hey yo!"))
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}

	// register client to a map of websocket clients
	clients[ws] = true
}

func main() {
	r := chi.NewMux()
	r.Use(middleware.Logger)
	r.Get("/", home)
	r.Get("/ws", wsHandler)

	s := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	log.Fatal(s.ListenAndServe())
}
