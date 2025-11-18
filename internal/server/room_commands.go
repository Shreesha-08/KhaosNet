package server

import (
	"fmt"
	"strings"
)

func RegisterRoomCommands(ch *CommandHandler) {
	ch.Register("/list", listUsersCommand)
	ch.Register("/name", changeNameCommand)
	ch.Register("/leave", leaveRoomCommand)
	ch.Register("/msg", privateMessageCommand)
}

func listUsersCommand(c *Client, args []string) {
	c.writeCh <- "Connected Users:"
	for i, _ := range c.currentRoom.broadcaster.clients {
		c.writeCh <- i
	}
}

func changeNameCommand(c *Client, args []string) {
	if len(args) == 1 {
		if c.CheckUniqueName(args[0]) {
			newName := args[0]
			oldName := c.name
			delete(c.currentRoom.broadcaster.clients, oldName)
			c.name = newName
			c.currentRoom.broadcaster.clients[newName] = c
			c.writeCh <- fmt.Sprintf("Username changed to %s", newName)
			c.currentRoom.broadcaster.msgCh <- &ClientMessage{msg: fmt.Sprintf("%s is now known as %s", oldName, newName), name: newName}
		} else {
			c.writeCh <- "Name already taken"
		}
	} else if len(args) == 1 {
		c.writeCh <- c.name
	} else {
		c.writeCh <- "Usage: /name or /name <newname>"
	}
}

func leaveRoomCommand(c *Client, args []string) {
	c.currentRoom.broadcaster.leaveCh <- c
	c.conn.Close()
}

func privateMessageCommand(c *Client, args []string) {
	if len(args) < 2 {
		c.writeCh <- "Usage: /msg <user> <message>"
		return
	}

	targetName := args[0]
	privateMsg := strings.Join(args[1:], " ")

	c.currentRoom.broadcaster.privateMsgCh <- &ClientMessage{
		msg:  fmt.Sprintf("(Private) %s: %s", c.name, privateMsg),
		name: targetName,
	}
}
