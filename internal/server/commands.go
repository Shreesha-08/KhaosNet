package server

import (
	"fmt"
	"strings"
)

type CommandFunc func(client *Client, args []string)

type CommandHandler struct {
	commands map[string]CommandFunc
}

func NewCommandHandler() *CommandHandler {
	cmdHandler := &CommandHandler{commands: make(map[string]CommandFunc)}
	cmdHandler.registerDefault()
	return cmdHandler
}

func (ch *CommandHandler) registerDefault() {
	ch.register("/list", listCommand)
	ch.register("/name", nameCommand)
	ch.register("/leave", leaveCommand)
	ch.register("/msg", pvtMsgCommand)
}

func (ch *CommandHandler) HandleCommand(c *Client, command string) bool {
	parts := strings.Fields(command)
	if len(parts) == 0 {
		return false
	}
	cmdName := parts[0]
	args := parts[1:]

	handler, exists := ch.commands[cmdName]
	if !exists {
		c.conn.Write([]byte("Unknown command. \n"))
		return true
	}

	handler(c, args)
	return true
}

func (ch *CommandHandler) register(name string, handler CommandFunc) {
	ch.commands[name] = handler
}

func listCommand(c *Client, args []string) {
	c.writeCh <- "Connected Users:"
	for i, _ := range c.broadcaster.clients {
		c.writeCh <- i
	}
}

func nameCommand(c *Client, args []string) {
	if len(args) == 1 {
		if c.CheckUniqueName(args[0]) {
			newName := args[0]
			oldName := c.name
			delete(c.broadcaster.clients, oldName)
			c.name = newName
			c.broadcaster.clients[newName] = c
			c.writeCh <- fmt.Sprintf("Username changed to %s", newName)
			c.broadcaster.msgCh <- &ClientMessage{msg: fmt.Sprintf("%s is now known as %s", oldName, newName), name: newName}
		} else {
			c.writeCh <- "Name already taken"
		}
	} else if len(args) == 1 {
		c.writeCh <- c.name
	} else {
		c.writeCh <- "Usage: /name or /name <newname>"
	}
}

func leaveCommand(c *Client, args []string) {
	c.broadcaster.leaveCh <- c
	c.conn.Close()
}

func pvtMsgCommand(c *Client, args []string) {
	if len(args) < 2 {
		c.writeCh <- "Usage: /msg <user> <message>"
		return
	}

	targetName := args[0]
	privateMsg := strings.Join(args[1:], " ")

	c.broadcaster.privateMsgCh <- &ClientMessage{
		msg:  fmt.Sprintf("(Private) %s: %s", c.name, privateMsg),
		name: targetName,
	}
}
