package server

import (
	"fmt"
	"strings"
	"time"
)

func RegisterRoomCommands(ch *CommandHandler) {
	ch.Register("/list", listUsersCommand)
	ch.Register("/name", changeNameCommand)
	ch.Register("/leave", leaveRoomCommand)
	ch.Register("/msg", privateMessageCommand)
	ch.Register("/kick", kickUser)
}

func listUsersCommand(c *Client, args []string) {
	users := make([]string, 0)

	for username := range c.currentRoom.broadcaster.clients {
		users = append(users, username)
	}

	resp := OutgoingMessage{
		Type: "users_list",
		Room: c.currentRoom.name,
		Data: map[string]any{
			"users": users,
		},
		Timestamp: time.Now().Unix(),
	}

	c.writeCh <- resp
}

func changeNameCommand(c *Client, args []string) {
	if len(args) == 0 {
		out := NewOutgoing(
			"system",
			"server",
			c.currentRoom.name,
			fmt.Sprintf("Your current name is: %s", c.name),
		)
		c.writeCh <- out
		return
	}

	if len(args) == 1 {
		newName := args[0]
		if !c.CheckUniqueName(newName) {
			out := NewOutgoing(
				"error",
				"server",
				c.currentRoom.name,
				"Name already taken",
			)
			c.writeCh <- out
			return
		}

		oldName := c.name
		delete(c.currentRoom.broadcaster.clients, oldName)
		c.name = newName
		c.currentRoom.broadcaster.clients[newName] = c

		ack := NewOutgoing(
			"system",
			"server",
			c.currentRoom.name,
			fmt.Sprintf("Username changed to %s", newName),
		)
		c.writeCh <- ack

		rename := NewOutgoing(
			"user_renamed",
			newName,
			c.currentRoom.name,
			fmt.Sprintf("%s is now known as %s", oldName, newName),
		)
		c.currentRoom.broadcaster.msgCh <- &ClientMessage{msg: rename, name: newName}
		return
	}

	out := NewOutgoing(
		"error",
		"server",
		c.currentRoom.name,
		"Usage: /name or /name <newname>",
	)
	c.writeCh <- out
}

func leaveRoomCommand(c *Client, args []string) {
	c.currentRoom.broadcaster.leaveCh <- c
	c.conn.Close()
}

func privateMessageCommand(c *Client, args []string) {
	if len(args) < 2 {
		out := NewOutgoing(
			"error",
			"server",
			c.currentRoom.name,
			"Usage: /msg <user> <message>",
		)
		c.writeCh <- out
		return
	}

	targetName := args[0]
	privateMsg := strings.Join(args[1:], " ")

	out := NewOutgoing("private", c.name, c.currentRoom.name, privateMsg)
	c.currentRoom.broadcaster.privateMsgCh <- &ClientMessage{
		msg:  out,
		name: targetName,
	}
}

func kickUser(c *Client, args []string) {
	if len(args) == 1 {
		if c.currentRoom.owner.name == c.name {
			targetName := args[0]
			if targetName == c.name {
				out := NewOutgoing(
					"error",
					"server",
					c.currentRoom.name,
					"Can't kick yourself out mate!",
				)
				c.writeCh <- out
				return
			}
			removedClient, exists := c.currentRoom.broadcaster.clients[targetName]
			if exists {
				c.currentRoom.broadcaster.leaveCh <- removedClient
			}
		} else {
			out := NewOutgoing(
				"error",
				"server",
				c.currentRoom.name,
				"Only room owner can kick users.",
			)
			c.writeCh <- out
		}
	}
}
