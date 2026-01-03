package main

import (
	"net/http"

	"github.com/gorilla/websocket"
)

var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan []byte)

func main() {
	upgrader := websocket.Upgrader{CheckOrigin: func(*http.Request) bool { return true }}

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, _ := upgrader.Upgrade(w, r, nil)
		clients[conn] = true

		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				delete(clients, conn)
				return
			}
			broadcast <- msg
		}
	})

	go func() {
		for msg := range broadcast {
			for c := range clients {
				c.WriteMessage(websocket.TextMessage, msg)
			}
		}
	}()

	http.ListenAndServe(":8080", nil)
}