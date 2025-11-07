package server

type Broadcaster struct {
	clients map[int]*Client
	joinCh  chan *Client
	leaveCh chan *Client
	msgCh   chan string
}

func NewBroadcaster() *Broadcaster {
	clientMap := make(map[int]*Client)
	join := make(chan *Client)
	leave := make(chan *Client)
	msg := make(chan string)
	return &Broadcaster{clients: clientMap, joinCh: join, leaveCh: leave, msgCh: msg}
}
