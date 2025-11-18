package server

import (
	"errors"
	"sync"
)

type Room struct {
	name        string
	broadcaster *Broadcaster
}

type RoomManager struct {
	rooms map[string]*Room
	mu    sync.Mutex
}

func NewRoomManager() *RoomManager {
	return &RoomManager{rooms: make(map[string]*Room)}
}

func (rm *RoomManager) CreateRoom(roomName string) (*Room, error) {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	if _, exists := rm.rooms[roomName]; exists {
		return nil, errors.New("room name already taken!")
	}
	newRoom := &Room{name: roomName, broadcaster: NewBroadcaster()}
	rm.rooms[roomName] = newRoom
	newRoom.broadcaster.Run()
	return newRoom, nil
}

func (rm *RoomManager) GetRoom(name string) *Room {
	rm.mu.Lock()
	defer rm.mu.Unlock()
	found, exists := rm.rooms[name]
	if exists {
		return found
	}
	return nil
}

func (rm *RoomManager) ListRooms() []string {
	rm.mu.Lock()
	defer rm.mu.Unlock()
	var rooms []string
	for room, _ := range rm.rooms {
		rooms = append(rooms, room)
	}
	return rooms
}
