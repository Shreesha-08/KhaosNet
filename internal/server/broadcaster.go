package server

import "fmt"

type Broadcaster struct {
	clients map[string]*Client
	joinCh  chan *Client
	leaveCh chan *Client
	msgCh   chan *ClientMessage
}

func NewBroadcaster() *Broadcaster {
	clientMap := make(map[string]*Client)
	join := make(chan *Client, 10)
	leave := make(chan *Client)
	msg := make(chan *ClientMessage, 10)
	return &Broadcaster{clients: clientMap, joinCh: join, leaveCh: leave, msgCh: msg}
}

func (b *Broadcaster) Run() {
	for {
		select {
		case client := <-b.joinCh:
			fmt.Println("New client joined")
			b.clients[client.name] = client
			client.writeCh <- fmt.Sprintf("Welcome, %s", client.name)
			msg := &ClientMessage{msg: fmt.Sprintf("%s joined!", client.name), name: client.name}
			b.msgCh <- msg
		case client := <-b.leaveCh:
			delete(b.clients, client.name)
			// close(client.writeCh)
			msg := &ClientMessage{msg: fmt.Sprintf("%s left!", client.name), name: client.name}
			b.msgCh <- msg
		case msg := <-b.msgCh:
			for client := range b.clients {
				if msg.name == b.clients[client].name {
					continue
				}
				b.clients[client].writeCh <- msg.msg
			}
		}
	}
}
