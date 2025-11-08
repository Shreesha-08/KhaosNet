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
	broadcaster *Broadcaster
}

func (c *Client) Read() {
	defer func() {
		c.broadcaster.leaveCh <- c
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
	c.broadcaster.joinCh <- c
	for {
		buf := make([]byte, 256)
		n, err := c.conn.Read(buf)
		if err != nil {
			return
		}
		msgStr := strings.TrimSpace(string(buf[:n]))
		if msgStr == "leave" {
			return
		} else {
			fmt.Println("Received msg: ", msgStr)
			msg := &ClientMessage{msg: fmt.Sprintf("%s: %s", c.name, msgStr), name: c.name}
			c.broadcaster.msgCh <- msg
		}
	}
}

func (c *Client) Write() {
	for msg := range c.writeCh {
		_, err := c.conn.Write([]byte(msg + "\n"))
		if err != nil {
			fmt.Println("Error writing to client:", c.name, err)
			c.broadcaster.leaveCh <- c
			return
		}
	}
}

func (c *Client) CheckUniqueName(name string) bool {
	_, exists := c.broadcaster.clients[name]
	return !exists
}
