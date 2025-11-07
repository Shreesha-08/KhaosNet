package server

import (
	"fmt"
	"net"
)

type Server struct {
	listener    net.Listener
	broadcaster *Broadcaster
	nextID      int
}

func (s *Server) NewServer(lis net.Listener) {
	s.listener = lis
	s.broadcaster = NewBroadcaster()
}

func (s *Server) AcceptConnection() {
	for {
		fmt.Println("New connection received")
		conn, err1 := s.listener.Accept()
		if err1 != nil {
			fmt.Println("Error getting connection")
		}
		clientName := fmt.Sprintf("User%d", s.nextID)
		s.nextID++
		client := &Client{conn: conn, name: clientName}
		go func() {
			s.broadcaster.joinCh <- client
		}()
	}
}
