package server

import (
	"fmt"
	"net"
	"strings"
)

type ClientMessage struct {
	msg  string
	name string
}

type Client struct {
	conn    net.Conn
	name    string
	writeCh chan string
	// doneCh      chan struct{}
	currentRoom *Room
	state       string
	server      *Server
}

func (c *Client) Read() {
	defer func() {
		c.currentRoom.broadcaster.leaveCh <- c
		close(c.writeCh)
		c.conn.Close()
	}()
	for {
		c.conn.Write([]byte("Enter your Username: "))
		namebuf := make([]byte, 64)
		n, err := c.conn.Read(namebuf)
		if err != nil {
			return
		}
		name := strings.TrimSpace(string(namebuf[:n]))
		if c.CheckUniqueName(name) {
			c.name = name
			break
		}
		c.conn.Write([]byte("Username already taken!\n"))
	}
	// c.broadcaster.joinCh <- c
	commandHandler := NewCommandHandler()
	for {
		buf := make([]byte, 256)
		n, err := c.conn.Read(buf)
		if err != nil {
			return
		}
		msgStr := strings.TrimSpace(string(buf[:n]))
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
		_, err := c.conn.Write([]byte(msg + "\n"))
		if err != nil {
			fmt.Println("Error writing to client:", c.name, err)
			c.currentRoom.broadcaster.leaveCh <- c
			return
		}
	}
}

func (c *Client) CheckUniqueName(name string) bool {
	// TODO: make this server level
	// _, exists := c.currentRoom.broadcaster.clients[name]
	// return !exists
	return true
}
