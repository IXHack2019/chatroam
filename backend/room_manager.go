package main

import (
	"fmt"
	//"log" //TODO: find extension to auto remove unused imports
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/asim/quadtree"

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
		if (nearClient.DeviceId != client.DeviceId) {
			nearRoom := nearClient.room

			if (len(nearRoom.members) < maxGroupSize) {
				client.room = nearRoom
				nearRoom.members = append(nearRoom.members, client)
			}
		}
	}

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
	// mutex.Lock()
	// defer func() {
	// 	mutex.Unlock()
	// }()

	for _, room := range rooms {
		curr := time.Now().UnixNano() / int64(time.Millisecond)
		if curr > room.expiry {
			//room is expired - reset it
			for _, clientInRoom := range room.members {
				freeClient(clientInRoom)

				// get a new room for the client
				getRoomForClient(clientInRoom)
			}
			room.expiry = time.Now().UnixNano()/int64(time.Millisecond) + 1000*20*1
		}
	}

	printRooms()
}

//for testing purposes - print status of each room
func printRooms() {

	fmt.Println("-------")
	for i, room := range rooms {
		fmt.Printf("%d people in room %d\n", len(room.members), i)
	}
	fmt.Println("-------")

}
