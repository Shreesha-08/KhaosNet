package server

import (
	"fmt"
	"strings"
)

type ClientMessage struct {
	msg  OutgoingMessage
	name string
}

type Client struct {
	conn        Connection
	name        string
	writeCh     chan OutgoingMessage
	usernameSet bool
	currentRoom *Room
	state       string
	server      *Server
}

func NewClient(conn Connection, s *Server, name string) *Client {
	return &Client{conn: conn, name: name, currentRoom: nil, state: "lobby", writeCh: make(chan OutgoingMessage, 10), server: s}
}

func (c *Client) Read() {
	defer c.Close()
	commandHandler := NewCommandHandler()
	for {
		incoming, err := c.conn.ReadAndGetData()
		if err != nil {
			return
		}
		if !c.usernameSet {
			if incoming.Command != "" && incoming.Command == "/username" {
				c.handleSetUsername(incoming.Args[0])
			} else {
				c.sendError("Please set your username first.")
			}
			continue
		}
		if incoming.Command != "" {
			commandHandler.HandleCommand(c, incoming.Command, incoming.Args)
			continue
		}

		if c.state == "lobby" {
			c.writeCh <- NewOutgoing("message", c.name, "lobby", "You need to join a chat room to send messages.")
			continue
		}
		out := NewOutgoing("message", c.name, c.currentRoom.name, incoming.Text)
		c.currentRoom.broadcaster.msgCh <- &ClientMessage{msg: out, name: c.name}
	}
}

func (c *Client) Write() {
	for msg := range c.writeCh {
		err := c.conn.Write(msg)
		if err != nil {
			fmt.Println("Error writing to client:", c.name, err)
			return
		}
	}
}

func (c *Client) Close() {
	if c.currentRoom != nil {
		c.currentRoom.broadcaster.leaveCh <- c
	}
	// close(c.writeCh)
	c.server.UnregisterClient(c)
	c.conn.Close()
}

func (c *Client) handleSetUsername(name string) {
	name = strings.TrimSpace(name)
	if name == "" {
		c.sendError("Username cannot be empty.")
		return
	}

	if !c.CheckUniqueName(name) {
		out := NewOutgoing("username_rejected", "server", "lobby", "Username already taken.")
		c.writeCh <- out
		return
	}
	c.name = name
	c.usernameSet = true
	c.server.RegisterClient(c)
	out := NewOutgoing("username_accepted", "server", "lobby", name)
	c.writeCh <- out

}

func (c *Client) sendError(text string) {
	out := NewOutgoing(
		"error",
		"server",
		"lobby",
		text,
	)
	c.writeCh <- out
}

func (c *Client) CheckUniqueName(name string) bool {
	_, exists := c.server.clients[name]
	return !exists
}
