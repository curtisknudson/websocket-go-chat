package main

import (
	"log"

	"net/http"

	"github.com/gorilla/websocket"
)

type Message struct {
	Username string `json:"username"`

	Text string `json:"text"`
}

func (m Message) Broadcast() {

	for client := range clients {

		err := client.WriteJSON(m)

		if err != nil {

			log.Println("Failed to write message to WebSocket:", err)

			client.Close()

			delete(clients, client)

		}

	}

}

var (

	// Utility from web socket library that upgrades the HTTP connection to 101 Switching Protocols code.

	upgrader = websocket.Upgrader{}

	// A way to reference the open clients with a connection/boolean key pair. If true the connection is open and messages will broadcast.

	clients = make(map[*websocket.Conn]bool)
)

func serveHTML(w http.ResponseWriter, r *http.Request) {

	http.ServeFile(w, r, "index.html")

	// http.ServeFile(w, r, "socketio.html")

}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {

	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {

		log.Println("Failed to upgrade the connection to WebSocket:", err)

		return

	}

	defer conn.Close()

	// Register new client

	clients[conn] = true

	// Continuously reads messages from the WebSocket connection and broadcasts them to all clients.

	for {

		var message Message

		err := conn.ReadJSON(&message)

		if err != nil {

			log.Println("Failed to read message from WebSocket:", err)

			delete(clients, conn)

			break

		}

		// Broadcast the message to all clients

		message.Broadcast()

	}

}

func main() {

	http.HandleFunc("/", serveHTML)

	http.HandleFunc("/ws", handleWebSocket)

	log.Println("Server started on http://localhost:8080")

	log.Fatal(http.ListenAndServe(":8080", nil))

}
