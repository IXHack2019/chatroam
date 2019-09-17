package main

import (
        "log"
        "net/http"
		"encoding/json"
		"log"
		"math/rand"

        "github.com/gorilla/websocket"
)

type Room struct {
	memberSockets  []*websocket.Conn
}
var roomAssignments = make(map[string]string)

var rooms []*Room

// Configure the upgrader
var upgrader = websocket.Upgrader{}


const (
	TypeConnect = iota
	TypeSend
	TypeReceive 
)

const maxGroupSize = 3
// Define our message object
type Message struct {
	Type int          `json:"type"`
	Data json.RawMessage `json:"data"`
}
type Connect struct {
	DeviceId string `json:"deviceId"`
}

type Send struct {
	Msg string `json:"msg"`
}

type Receive struct {
	Msg string `json:"msg"`
	Username string `json:"username"`
}

type User struct {
	DeviceId string
	Username string
}

type Client struct {
	socket *websocket.Conn
	room *Room
}

func main() {
	log.SetFlags(log.LstdFlags)
	http.HandleFunc("/connect", handleMessage)
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}

func handleMessage(w http.ResponseWriter, r *http.Request) {
	// Upgrade initial GET request to a websocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	// Make sure we close the connection when the function returns
	defer ws.Close()

	for {

		_, request, err := ws.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			return
		}

		var message Message
		err = json.Unmarshal(request, &message)
		if err != nil {
			log.Printf("Error unmarshalling message: %s", err)
			continue
		}

		log.Printf("Parsed message: %s", message)

		if message.Type == TypeConnect {

			handleConnect(message.Data)

		} else if message.Type == TypeSend {

		} else if message.Type == TypeReceive {


		}


	}

}

func (client *Client) handleConnect(data json.RawMessage) {
	var connect Connect
	err := json.Unmarshal(data, &connect)
	if err != nil {
		log.Printf("Error unmarshalling host connect: %s", err)
		return
	}

	numRooms := len(rooms)
	if numRooms == 0 {
		rooms = append(rooms, &Room{
			participants: []*websocket.Conn{client.socket},
		})
	}
	for i, room := range rooms {
		if (len(room.participants) < maxGroupSize) { //TODO race condition here lul
			room.memberSocket = append(room.memberSocket, client.socket)
		}
	} 
	connect.DeviceId

}
func handleConnections(w http.ResponseWriter, r *http.Request) {
	// Upgrade initial GET request to a websocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	// Make sure we close the connection when the function returns
	defer ws.Close()

	// Register our new client
	clients[ws] = true
	for {
		var msg Message
		// Read in a new message as JSON and map it to a Message object
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("error: %v", err)
			delete(clients, ws)
			break
		}
		// Send the newly received message to the broadcast channel
		broadcast <- msg
	}
}

func handleMessages() {
	for {
		// Grab the next message from the broadcast channel
		msg := <-broadcast
		// Send it out to every client that is currently connected
		for client := range clients {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Printf("error: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}
