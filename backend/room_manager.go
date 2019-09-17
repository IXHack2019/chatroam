package main

import (
	"sync"
)

var pool = &sync.Pool{
	// New creates an object when the pool has nothing available to return.
	New: func() interface{} {
		return &Room{
			members: []*Client{},
		}
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
