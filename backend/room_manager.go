package main

import (
	"fmt"
	"log"
	"math"
	"sync"
	"time"

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

	minDistance := math.MaxFloat64
	var minRoom *Room = nil

	for i, room := range rooms {
		curr := time.Now().UnixNano() / int64(time.Millisecond)
		if len(room.members) < maxGroupSize && curr <= room.expiry {
			if len(room.members) != 0 {
				firstMember := room.members[0]

				distance := distanceInKmBetweenEarthCoordinates(firstMember.Lat, firstMember.Lon, client.Lat, client.Lon)

				log.Printf("Room %d lat1 %f lon1 %f lat2 %f lon2 %f Distance: %f\n", i, firstMember.Lat, firstMember.Lon, client.Lat, client.Lon, distance)

				if distance < minDistance || minDistance < 0 {
					minDistance = distance
					minRoom = room
				}
			} else if minRoom == nil {
				//room is empty, use it as last resort
				//treat -1 as infinity
				minDistance = -1
				minRoom = room
			}
		}
	}

	if minRoom != nil {
		minRoom.members = append(minRoom.members, client)
		client.room = minRoom
	}

	//no room was available - create new room
	if client.room == nil {
		newRoom := &Room{
			members: []*Client{},
			expiry:  time.Now().UnixNano()/int64(time.Millisecond) + 1000*20*1,
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
		}
	}

	printRooms()

	return success
}

// go through all rooms, vacate the expired rooms and put the clients into new rooms
func resetRooms() {
	mutex.Lock()
	defer func() {
		mutex.Unlock()
	}()

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
