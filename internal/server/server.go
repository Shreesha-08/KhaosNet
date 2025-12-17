package server

import (
	"fmt"
	"net"
	"net/http"
	"sync"

	"github.com/coder/websocket"
)

type Server struct {
	listener net.Listener
	roomMgr  *RoomManager
	nextID   int
	mu       sync.Mutex
	clients  map[string]*Client
}

func (s *Server) NewServer(lis net.Listener) {
	s.listener = lis
	s.roomMgr = NewRoomManager()
	s.clients = make(map[string]*Client)
}

func (s *Server) Start() {
	go s.AcceptConnections() // TCP + WebSocket
	http.HandleFunc("/ws", s.WebSocketHandler)
	fmt.Println("WebSocket server on :8080")
	http.ListenAndServe(":8080", nil)
}

func (s *Server) AcceptConnections() {
	// for {
	// 	conn, err1 := s.listener.Accept()
	// 	if err1 != nil {
	// 		fmt.Println("Error getting connection")
	// 	}
	// 	clientName := s.GenerateName()
	// 	tcpConn := NewTCPConnection(conn)
	// 	client := NewClient(tcpConn, s, clientName)
	// 	go client.Read()
	// 	go client.Write()
	// }
}

func (s *Server) WebSocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		InsecureSkipVerify: true, // okay for local development
	})
	if err != nil {
		return
	}

	ws := NewWebSocketConn(conn)

	clientName := s.GenerateName()
	client := NewClient(ws, s, clientName)
	go client.Read()
	go client.Write()
}

func (s *Server) GenerateName() string {
	name := fmt.Sprintf("User%d", s.nextID)
	s.nextID++
	return name
}

func (s *Server) RegisterClient(c *Client) {
	s.mu.Lock()
	s.clients[c.name] = c
	s.mu.Unlock()
}

func (s *Server) UnregisterClient(c *Client) {
	s.mu.Lock()
	delete(s.clients, c.name)
	s.mu.Unlock()
}

func (s *Server) Broadcast(msg OutgoingMessage) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, cl := range s.clients {
		select {
		case cl.writeCh <- msg:
		default:
		}
	}
}
