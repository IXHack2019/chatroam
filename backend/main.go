package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	// "log"
	// "math/rand"

	"github.com/gorilla/websocket"
)

type Room struct {
	members []*Client
	expiry  int64
}

// Configure the upgrader
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

const (
	TypeConnect = iota
	TypeSend
	TypeQuery
)

const maxGroupSize = 3

// Define our message object
type Message struct {
	Type int             `json:"type"`
	Data json.RawMessage `json:"data"`
}
type Connect struct {
	DeviceId string  `json:"deviceId"`
	Lon      float64 `json:"lon"`
	Lat      float64 `json:"lat"`
}

type ReceivedMessage struct {
	DeviceId string `json:"deviceId"`
	Msg      string `json:"msg"`
	Username string `json:"username"` // Problem: User may be able to change this with inspector
}

type RegistrationResponse struct {
	Type     int
	Username string `json:"username"`
}

type QueryResponse struct {
	Records []ClientRecord `json:"records"`
}

type ClientRecord struct {
	Username string  `json:"username"`
	Lon      float64 `json:"lon"`
	Lat      float64 `json:"lat"`
	RoomID   int     `json:"roomID"`
	LastMsg  string  `json:"lastMsg"`
}

type Client struct {
	socket   *websocket.Conn
	room     *Room
	DeviceId string
	Username string
	Lon      float64
	Lat      float64
	LastMsg  string
}

var connectedClients = make(map[string]*Client)
var rooms []*Room

func main() {
	//TODO: fix reset not putting users in a new room
	//go scheduler(time.NewTicker(time.Second * 5))
	log.SetFlags(log.LstdFlags)
	http.HandleFunc("/connect", handleMessage)
	log.Fatal(http.ListenAndServe("localhost:8888", nil))
}

func scheduler(tick *time.Ticker) {
	for range tick.C {
		resetRooms()
	}
}

func handleMessage(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received message\n")

	// Upgrade initial GET request to a websocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	// Make sure we close the connection when the function returns
	defer ws.Close()

	var client *Client
	defer func() {
		if client != nil {
			log.Printf("Removing client %s from room", client.DeviceId)
			freeClient(client)
		}
	}()

	for {
		var message Message
		err = ws.ReadJSON(&message)

		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Error unmarshalling message: %s", err)
			}
			return
		}

		log.Printf("Parsed message: %s type: %d", message.Data, message.Type)

		if message.Type == TypeConnect {
			client = &Client{
				socket: ws,
			}

			client.handleConnect(message.Data)
		} else if message.Type == TypeSend {
			var receivedMessage ReceivedMessage
			err := json.Unmarshal(message.Data, &receivedMessage)
			if err != nil {
				log.Printf("Error unmarshalling host connect: %s", err)
				return
			}

			client.LastMsg = receivedMessage.Msg

			for _, member := range client.room.members {
				log.Printf("Writing to deviceId %s's socket: %s", member.DeviceId, receivedMessage.Msg)

				member.socket.WriteJSON(message)
			}

		} else if message.Type == TypeQuery {
			client.handleQuery()
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

	getRoomForClient(client)

	// minDistance := float64(0)
	// var minRoom *Room = nil
	// for _, room := range rooms { // big performance issue here if number of rooms is large, but this is a hackathon

	// 	if len(room.members) < maxGroupSize { // TODO race condition here lul

	// 		firstMember := room.members[0]

	// 		distance := distanceInKmBetweenEarthCoordinates(firstMember.Lat, firstMember.Lon, client.Lat, client.Lon)

	// 		log.Printf("Distance: %f", distance)

	// 		if distance < minDistance {
	// 			minDistance = distance
	// 			minRoom = room
	// 		}
	// 	}
	// }

	// if minRoom != nil {
	// 	minRoom.members = append(minRoom.members, client)
	// 	client.room = minRoom
	// }

	// if client.room == nil {
	// 	newRoom := &Room{
	// 		members: []*Client{client},
	// 	}
	// 	rooms = append(rooms, newRoom)
	// 	client.room = newRoom
	// }
	client.Lat = connect.Lat
	client.Lon = connect.Lon
	client.DeviceId = connect.DeviceId
	client.Username = getRandomName()
	connectedClients[connect.DeviceId] = client

	client.socket.WriteJSON(RegistrationResponse{0, client.Username})
}

func (client *Client) handleQuery() {
	response := QueryResponse{}

	for i, room := range rooms {
		for _, client := range room.members {
			response.Records = append(response.Records, ClientRecord{
				Username: client.Username,
				Lon:      client.Lon,
				Lat:      client.Lat,
				RoomID:   i,
				LastMsg:  client.LastMsg,
			})
		}
	}
}
