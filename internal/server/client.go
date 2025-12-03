package server

import (
	"fmt"
	"strings"
)

type ClientMessage struct {
	msg  string
	name string
}

type Client struct {
	conn    Connection
	name    string
	writeCh chan string
	// doneCh      chan struct{}
	currentRoom *Room
	state       string
	server      *Server
}

func NewClient(conn Connection, s *Server, name string) *Client {
	return &Client{conn: conn, name: name, currentRoom: nil, state: "lobby", writeCh: make(chan string, 10), server: s}
}

func (c *Client) Read() {
	defer func() {
		c.Close()
	}()
	for {
		c.conn.Write("Enter your Username: ")
		name, err := c.conn.Read()
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		if c.CheckUniqueName(name) {
			c.name = name
			break
		}
		c.conn.Write("Username already taken!\n")
	}
	// c.broadcaster.joinCh <- c
	commandHandler := NewCommandHandler()
	for {
		msgStr, err := c.conn.Read()
		if err != nil {
			return
		}
		if strings.HasPrefix(msgStr, "/") {
			commandHandler.HandleCommand(c, msgStr)
			continue
		}

		if c.state == "lobby" {
			c.writeCh <- "You need to join a chat room to send messages. "
			continue
		}
		fmt.Println("Received msg: ", msgStr)
		msg := &ClientMessage{msg: fmt.Sprintf("%s: %s", c.name, msgStr), name: c.name}
		c.currentRoom.broadcaster.msgCh <- msg
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
	c.conn.Close()
}

func (c *Client) CheckUniqueName(name string) bool {
	// TODO: make this server level
	// _, exists := c.currentRoom.broadcaster.clients[name]
	// return !exists
	return true
}
