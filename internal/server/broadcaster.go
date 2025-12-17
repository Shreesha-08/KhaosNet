package server

import (
	"fmt"
)

type Broadcaster struct {
	clients      map[string]*Client
	joinCh       chan *Client
	leaveCh      chan *Client
	msgCh        chan *ClientMessage
	privateMsgCh chan *ClientMessage
	roomName     string
}

func NewBroadcaster(roomName string) *Broadcaster {
	clientMap := make(map[string]*Client)
	join := make(chan *Client, 10)
	leave := make(chan *Client)
	msg := make(chan *ClientMessage, 10)
	privateMsg := make(chan *ClientMessage, 10)
	return &Broadcaster{clients: clientMap, joinCh: join, leaveCh: leave, msgCh: msg, privateMsgCh: privateMsg, roomName: roomName}
}

func (b *Broadcaster) Run() {
	for {
		select {
		case client := <-b.joinCh:
			b.clients[client.name] = client
			welcome := NewOutgoing("system", "server", b.roomName, fmt.Sprintf("Welcome, %s", client.name))
			client.writeCh <- welcome
			join := NewOutgoing("user_joined", client.name, b.roomName, fmt.Sprintf("%s joined!", client.name))
			b.msgCh <- &ClientMessage{msg: join, name: client.name}

		case client := <-b.leaveCh:
			delete(b.clients, client.name)
			leave := NewOutgoing("user_left", client.name, b.roomName, fmt.Sprintf("%s left!", client.name))
			b.msgCh <- &ClientMessage{msg: leave, name: client.name}
			msgForClient := NewOutgoing("left_room", client.name, b.roomName, fmt.Sprintf("Left %s!", b.roomName))
			client.writeCh <- msgForClient

		case msg := <-b.msgCh:
			for _, cl := range b.clients {
				if msg.name == cl.name {
					continue
				}
				select {
				case cl.writeCh <- msg.msg:
				default:
					fmt.Printf("Dropping message for %s (writeCh full)\n", cl.name)
				}
			}
		case privateMsg := <-b.privateMsgCh:
			target, exists := b.clients[privateMsg.name]
			if !exists {
				continue
			}

			select {
			case target.writeCh <- privateMsg.msg:
			default:
				fmt.Printf("Dropping private message for %s (writeCh full)\n", privateMsg.name)
			}
		}
	}
}
