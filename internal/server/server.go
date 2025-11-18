package server

import (
	"fmt"
	"net"
)

type Server struct {
	listener net.Listener
	roomMgr  *RoomManager
	nextID   int
}

func (s *Server) NewServer(lis net.Listener) {
	s.listener = lis
	s.roomMgr = NewRoomManager()
}

func (s *Server) Start() {
	s.AcceptConnections()
}

func (s *Server) AcceptConnections() {
	for {
		fmt.Println("New connection received")
		conn, err1 := s.listener.Accept()
		if err1 != nil {
			fmt.Println("Error getting connection")
		}
		clientName := fmt.Sprintf("User%d", s.nextID)
		s.nextID++
		client := &Client{conn: conn, name: clientName, currentRoom: nil, state: "lobby", writeCh: make(chan string, 10), server: s}
		go client.Read()
		go client.Write()
	}
}
