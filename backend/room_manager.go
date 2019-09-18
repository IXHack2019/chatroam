package main

import (
	"fmt"
	"sync"
	"time"
)

var mutex = &sync.Mutex{}

func getRoomForClient(client *Client) {

	for _, room := range rooms {
		mutex.Lock()
		if len(room.members) < maxGroupSize {
			client.room = room
			room.members = append(room.members, client)
		}
		mutex.Unlock()
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

	printRooms()

}

func freeClient(client *Client) bool {

	room := client.room
	client.room = nil
	success := false

	mutex.Lock()

	//delete the client from room members
	for i, clientInRoom := range room.members {
		if clientInRoom == client {
			room.members = append(room.members[:i], room.members[i+1:]...)
			success = true
		}
	}

	mutex.Unlock()

	printRooms()

	return success
}

// go through all rooms, vacate the expired rooms and put the clients into new rooms
func resetRooms() {

	for _, room := range rooms {
		mutex.Lock()
		curr := time.Now().UnixNano() / int64(time.Millisecond)
		if curr > room.expiry {
			//room is expired - reset it
			room.members = room.members[:0]
			for _, clientInRoom := range room.members {
				// get a new room for the client
				getRoomForClient(clientInRoom)
			}
			room.expiry = time.Now().UnixNano() / int64(time.Millisecond)
		}
		mutex.Unlock()
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
