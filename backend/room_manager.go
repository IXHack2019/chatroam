package main

import (
	"sync"
	"time"
)

var pool = &sync.Pool{
	// New creates an object when the pool has nothing available to return.
	New: func() interface{} {
		newRoom := &Room{
			members: []*Client{},
			expiry:  time.Now().UnixNano() / int64(time.Millisecond),
		}
		rooms = append(rooms, newRoom)
		return newRoom
	},
}

var mutex = &sync.Mutex{}

func getRoomForClient(client *Client) {

	room := pool.Get().(*Room)
	client.room = room
	room.members = append(room.members, client)

	//if room still has capacity, put it back in the pool
	if len(room.members) < maxGroupSize {
		pool.Put(room)
	}

}

func freeClient(client *Client) bool {

	room := client.room
	client.room = nil
	success := false

	mutex.Lock()
	oldLen := len(room.members)

	//delete the client from room members
	for i, clientInRoom := range room.members {
		if clientInRoom == client {
			room.members = append(room.members[:i], room.members[i+1:]...)
			success = true
		}
	}

	//if room now has capacity put it back in the pool
	if oldLen >= maxGroupSize && len(room.members) < maxGroupSize {
		pool.Put(room)
	}
	mutex.Unlock()

	return success
}

// go through all rooms, vacate the expired rooms and put the clients into new rooms
func resetRooms() {

	for _, room := range rooms {
		mutex.Lock()
		curr := time.Now().UnixNano() / int64(time.Millisecond)
		if curr > room.expiry {
			//room is expired - reset it
			for _, clientInRoom := range room.members {
				// get a new room for the client
				clientInRoom.room = nil
				getRoomForClient(clientInRoom)
			}
			room.members = room.members[:0]
			room.expiry = time.Now().UnixNano() / int64(time.Millisecond)
		}
		mutex.Unlock()
	}

}
