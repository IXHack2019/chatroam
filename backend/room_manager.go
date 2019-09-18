package main

import (
	"fmt"
	"log" //TODO: find extension to auto remove unused imports
	"sync"
	"time"

	"github.com/asim/quadtree"
	"github.com/gorilla/websocket"
)

//type RoomManager struct {
//	rooms: []*Room,
//	mutex: *sync.Mutex{}
//}
//
//func NewRoomManager() *RoomManager {
//	return &RoomManager{
//		mutex: sync.Mutex{}
//	}
//}

var mutex = &sync.Mutex{}

func getRoomForClient(client *Client) {
	mutex.Lock()
	defer func() {
		mutex.Unlock()
	}()

	center := quadtree.NewPoint(client.Lat, client.Lon, nil)
	bounds := quadtree.NewAABB(center, center.HalfPoint(maxSearchDistance))

	maxPoints := 10 // Try this many points before giving up
	for _, point := range qtree.KNearest(bounds, maxPoints, nil) {
		//log.Printf("Found point: %s\n", point.Data().(string))

		nearClient := point.Data().(*Client)
		if nearClient.DeviceId != client.DeviceId {
			nearRoom := nearClient.room

			if len(nearRoom.members) < maxGroupSize {
				client.room = nearRoom
				nearRoom.members = append(nearRoom.members, client)
			}
		}
	}
	// if client.room == nil && len(rooms) > 0 {
	// 	client.room = rooms[0]
	// 	rooms[0].members = append(rooms[0].members, client)
	// }
	rooms = nil // TODO: make 
	//no room was available - create new room
	if client.room == nil {
		newRoom := &Room{
			members: []*Client{},
			expiry:  time.Now().UnixNano()/int64(time.Millisecond) + 1000*20*1,
			mutex:   &sync.Mutex{},
		}
		newRoom.members = append(newRoom.members, client)
		client.room = newRoom
		rooms = append(rooms, newRoom)
	}

	for _, message := range client.room.messages {
		if client.socket != nil {
			client.socket.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf(`{"type": 1, "data": {"username": "%s", "msg":"%s" } }`, message.Name, message.Text)))
		}
	}

	printRooms()

}

func freeClient(client *Client) bool {
	mutex.Lock()
	defer func() {
		mutex.Unlock()
	}()

	room := client.room

	if room == nil {
		//client not in room - do nothing
		return true
	}

	//client.room = nil
	success := false

	//delete the client from room members
	for i, clientInRoom := range room.members {
		if clientInRoom == client {
			room.members = append(room.members[:i], room.members[i+1:]...)
			success = true
			break
		}
	}

	printRooms()

	return success
}

// go through all rooms, vacate the expired rooms and put the clients into new rooms
func resetRooms() {
	// mutex.Lock() // TODO: synbc this stuff
	// defer func() {
	// 	mutex.Unlock()
	// }()
	var freeClients = make(map[int]*Client)

	
	for _, room := range rooms {
		// curr := time.Now().UnixNano() / int64(time.Millisecond)
		// if curr > room.expiry {
		//room is expired - reset it
		for i, clientInRoom := range room.members {
			freeClients[i] = clientInRoom
			clientInRoom.room = nil
		}

		room.members = nil
		// room.expiry = time.Now().UnixNano()/int64(time.Millisecond) + 1000*20*1
		// }
	}

	for _, client := range freeClients {
		// get a new room for the client
		log.Println(client.Username);
		getRoomForClient(client)
	}

	printRooms()
}

//for testing purposes - print status of each room
func printRooms() {

	fmt.Println("-------")
	for i, room := range rooms {
		fmt.Printf("%d people in room %d\n", len(room.members), i)
		for _, client := range room.members {
			fmt.Printf("%s in room %d\n", client.Username, i)
		}
	}
	fmt.Println("-------")

}
