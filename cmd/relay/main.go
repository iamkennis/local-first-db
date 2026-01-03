package main

import (
	"net/http"

	"github.com/gorilla/websocket"
)

var peers = map[*websocket.Conn]bool{}
var up = websocket.Upgrader{CheckOrigin: func(*http.Request) bool { return true }}

func handler(w http.ResponseWriter, r *http.Request) {
	c, _ := up.Upgrade(w, r, nil)
	peers[c] = true

	for {
		_, msg, err := c.ReadMessage()
		if err != nil {
			delete(peers, c)
			return
		}
		for p := range peers {
			p.WriteMessage(websocket.TextMessage, msg)
		}
	}
}

func main() {
	http.HandleFunc("/ws", handler)
	http.ListenAndServe(":8081", nil)
}