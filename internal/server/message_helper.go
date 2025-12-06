package server

import "time"

func NewOutgoing(t, from, room, text string) OutgoingMessage {
	return OutgoingMessage{
		Type:      t,
		From:      from,
		Room:      room,
		Text:      text,
		Timestamp: time.Now().Unix(),
	}
}
