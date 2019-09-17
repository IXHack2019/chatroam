package main

import (
        "log"
        "net/http"
		"encoding/json"
		// "log"
		// "math/rand"

        "github.com/gorilla/websocket"
)

type Room struct {
	members  []*Client
}

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

type Client struct {
	socket *websocket.Conn 
	room *Room
	DeviceId string 
	Username string
}

type Response struct {
	Status int `json:"status"`
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
		var message Message
		err = ws.ReadJSON(&message)
		if err != nil {
			log.Printf("Error unmarshalling message: %s", err)
			continue
		}

		log.Printf("Parsed message: %s", message)

		if message.Type == TypeConnect {
			client := Client{
				socket: ws,
				
			}

			client.handleConnect(message.Data)

		} else if message.Type == TypeSend {

		} else if message.Type == TypeReceive {


		}

		websocket.WriteJSON(ws, Response{Status: 0})
	}


}

func (client *Client) handleConnect(data json.RawMessage) {
	var connect Connect
	err := json.Unmarshal(data, &connect)
	if err != nil {
		log.Printf("Error unmarshalling host connect: %s", err)
		return
	}

	for _, room := range rooms {
		if (len(room.members) < maxGroupSize) { //TODO race condition here lul
			room.members = append(room.members, client)
			client.room = room
			break
		}
	} 

	if client.room == nil {
		newRoom := &Room{
			members: []*Client{client},
		}
		rooms = append(rooms, newRoom)
		client.room = newRoom
	}
	
}
// func handleConnections(w http.ResponseWriter, r *http.Request) {
// 	// Upgrade initial GET request to a websocket
// 	ws, err := upgrader.Upgrade(w, r, nil)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	// Make sure we close the connection when the function returns
// 	defer ws.Close()

// 	// Register our new client
// 	clients[ws] = true
// 	for {
// 		var msg Message
// 		// Read in a new message as JSON and map it to a Message object
// 		err := ws.ReadJSON(&msg)
// 		if err != nil {
// 			log.Printf("error: %v", err)
// 			delete(clients, ws)
// 			break
// 		}
// 		// Send the newly received message to the broadcast channel
// 		broadcast <- msg
// 	}
// }

// func handleMessages() {
// 	for {
// 		// Grab the next message from the broadcast channel
// 		msg := <-broadcast
// 		// Send it out to every client that is currently connected
// 		for client := range clients {
// 			err := client.WriteJSON(msg)
// 			if err != nil {
// 				log.Printf("error: %v", err)
// 				client.Close()
// 				delete(clients, client)
// 			}
// 		}
// 	}
// }
