package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var (
	clients   = make(map[*websocket.Conn]bool) // connected clients
	broadcast = make(chan Log)                 // broadcast channel
	upgrader  = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	websocketTimeFormat = "15:04:05"
)

func handleConnections(w http.ResponseWriter, r *http.Request) {
	// Upgrade initial GET request to a websocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	// Make sure we close the connection when the function returns
	// defer ws.Close()

	// Register our new client
	clients[ws] = true
}

type websocketLogMessage struct {
	Time string `json:"time"`
	Hits int64  `json:"hits"`
	Code int    `json:"code"`
}

func handleMessages() {
	for {
		// Grab the next message from the broadcast channel
		msg := <-broadcast
		wsMsg := websocketLogMessage{
			Time: msg.Time.Format(websocketTimeFormat),
			Hits: msg.Hits,
			Code: msg.Code,
		}

		// Send it out to every client that is currently connected
		for client := range clients {
			err := client.WriteJSON(wsMsg)
			if err != nil {
				log.Printf("error: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}
