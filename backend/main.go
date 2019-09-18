package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"./generation"

	"github.com/gorilla/websocket"
	"github.com/asim/quadtree"
)

type Room struct {
	members  []*Client
	expiry   int64
	messages []RoomMessage
	mutex    *sync.Mutex
}

func (r *Room) AddMessage(message RoomMessage) {
	r.mutex.Lock()
	defer func() {
		r.mutex.Unlock()
	}()

	r.messages = append(r.messages, message)
	if len(r.messages) > 10 {
		r.messages = r.messages[1:]
	}
}

type RoomMessage struct {
	Name string
	Text string
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
	TypeRoom
)

const maxGroupSize = 3
const maxSearchDistance = 10000
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
	Type     int    `json:"type"`
	Username string `json:"username"`
}

type QueryResponse struct {
	Type    int            `json:"type"`
	Records []ClientRecord `json:"records"`
}

type ClientRecord struct {
	Name    string  `json:"name"`
	Lon     float64 `json:"lon"`
	Lat     float64 `json:"lat"`
	RoomID  int     `json:"roomID"`
	LastMsg string  `json:"lastMsg"`
}

type RoomResponse struct {
	Type    int          `json:"type"`
	Records []RoomRecord `json:"records"`
}

type RoomRecord struct {
	ID        int          `json:"roomID"`
	Log       []string     `json:"chatlog"`
	Locations [][2]float64 `json:"locations"`
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

var rooms []*Room
var qtree *quadtree.QuadTree
var connectedClients = make(map[string]*Client)
var botClients = []*Client{
	&Client{
		Username: "Matt Lewis",
		Lat:      43.652375,
		Lon:      -79.376576,
		LastMsg:  "This app is siiiick! YEET!",
	},
	&Client{
		Username: "Charles Black",
		Lat:      43.652238,
		Lon:      -79.380588,
		LastMsg:  "Go Redskins!!!",
	},
	&Client{
		Username: "Frank Castle",
		Lat:      43.652112,
		Lon:      -79.380688,
		LastMsg:  "It's time for punishment",
	},
	&Client{
		Username: "George Foreman",
		Lat:      43.652438,
		Lon:      -79.380388,
		LastMsg:  "Cook a steak in 5 minutes!",
	},
	&Client{
		Username: "Gloria Raynor",
		Lat:      43.651238,
		Lon:      -79.381588,
		LastMsg:  "Respect!",
	},
	&Client{
		Username: "Marie Curie",
		Lat:      43.651438,
		Lon:      -79.381288,
		LastMsg:  "Yeah! Science!",
	},
	&Client{
		Username: "Johan Strutt",
		Lat:      43.653238,
		Lon:      -79.384588,
		LastMsg:  "Howdy There",
	},
	&Client{
		Username: "Julio Jones",
		Lat:      43.653238,
		Lon:      -79.382588,
		LastMsg:  "Matty Ice hit me up in the end zone!",
	},
	&Client{
		Username: "Adrian Peterson",
		Lat:      43.652510,
		Lon:      -79.390468,
		LastMsg:  "Wow. An orange peanut? For me? Wow",
	},
	&Client{
		Username: "Tedd George",
		Lat:      43.649161,
		Lon:      -79.375986,
		LastMsg:  "Anyone interested in grabbing some dim sum?",
	},
	&Client{
		Username: "Lucy Diamond",
		Lat:      43.655029,
		Lon:      -79.370690,
		LastMsg:  "Looking for someone to go shopping with!!!",
	},
}

func main() {
	centerPoint := quadtree.NewPoint(0.0, 0.0, nil)
	halfPoint := quadtree.NewPoint(90.0, 180.0, nil)
	boundingBox := quadtree.NewAABB(centerPoint, halfPoint)
	
	qtree = quadtree.New(boundingBox, 0, nil)
	// for _, bot := range botClients {
	// 	getRoomForClient(bot)
	// }

	//go scheduler(time.NewTicker(time.Second * 5))
	//go resetScheduler(time.NewTicker(time.Second * 5))

	log.SetFlags(log.LstdFlags)
	http.HandleFunc("/connect", handleMessage)
	log.Fatal(http.ListenAndServe("localhost:8888", nil))
}

func scheduler(tick *time.Ticker) {
	for range tick.C {
		updateTestClients()
	}
}

func resetScheduler(tick *time.Ticker) {
	for range tick.C {
		resetRooms()
	}
}

func updateTestClients() {
	for _, client := range connectedClients {
		offset := float64(rand.Intn(100))*0.00001 - 0.0005
		client.Lat += offset
		offset = float64(rand.Intn(100))*0.00001 - 0.0005
		client.Lon += offset

		time.Sleep(200 * time.Millisecond)
	}

	for _, testClient := range botClients {
		testClient.LastMsg = generation.GetRandomMessage()
		for _, client := range testClient.room.members {
			if client.socket != nil {
				client.socket.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf(`{"type": 1, "data": {"username": "%s", "msg":"%s" } }`, testClient.Username, testClient.LastMsg)))
			}
		}
		testClient.room.AddMessage(RoomMessage{
			Name: testClient.Username,
			Text: testClient.LastMsg,
		})

		offset := float64(rand.Intn(100))*0.00001 - 0.0005
		testClient.Lat += offset
		offset = float64(rand.Intn(100))*0.00001 - 0.0005
		testClient.Lon += offset

		time.Sleep(200 * time.Millisecond)
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

	var client = &Client{
		socket: ws,
	}
	defer func() {
		log.Printf("Removing client %s from room", client.DeviceId)
		freeClient(client)
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
			client.handleConnect(message.Data)
		} else if message.Type == TypeSend {
			client.handleSend(message)
		} else if message.Type == TypeQuery {
			client.handleQuery()
		} else if message.Type == TypeRoom {
			client.handleRoom()
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
	client.Lat = connect.Lat
	client.Lon = connect.Lon
	client.DeviceId = connect.DeviceId
	client.Username = generation.GetRandomName()

	point := quadtree.NewPoint(client.Lat, client.Lon, client)
	if !qtree.Insert(point) {
		log.Fatal("Failed to insert the point")
	}

	connectedClients[connect.DeviceId] = client

	client.socket.WriteJSON(RegistrationResponse{0, client.Username})

	getRoomForClient(client)
}

func (client *Client) handleSend(message Message) {
	var receivedMessage ReceivedMessage
	err := json.Unmarshal(message.Data, &receivedMessage)
	if err != nil {
		log.Printf("Error unmarshalling host connect: %s", err)
		return
	}

	client.LastMsg = receivedMessage.Msg
	if client.room != nil {
		client.room.AddMessage(RoomMessage{
			Name: client.Username,
			Text: receivedMessage.Msg,
		})
	}

	for _, member := range client.room.members {
		if member != nil && member.socket != nil {
			log.Printf("Writing to deviceId %s's socket: %s", member.DeviceId, receivedMessage.Msg)

			member.socket.WriteJSON(message)
		}
	}
}

func (client *Client) handleQuery() {
	response := QueryResponse{Type: TypeQuery}

	for i, room := range rooms {
		for _, client := range room.members {
			response.Records = append(response.Records, ClientRecord{
				Name:    client.Username,
				Lon:     client.Lon,
				Lat:     client.Lat,
				RoomID:  i,
				LastMsg: client.LastMsg,
			})
		}
	}

	client.socket.WriteJSON(response)
}

func (client *Client) handleRoom() {
	response := RoomResponse{Type: TypeRoom}

	for i, room := range rooms {
		record := RoomRecord{ID: i}

		var log []string
		for _, message := range room.messages {
			log = append(log, fmt.Sprintf("%s: %s", message.Name, message.Text))
		}
		record.Log = log

		for _, client := range room.members {
			record.Locations = append(record.Locations, [2]float64{client.Lat, client.Lon})
		}

		response.Records = append(response.Records, record)
	}

	client.socket.WriteJSON(response)
}
