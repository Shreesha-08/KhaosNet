package server

import (
	"fmt"
)

type CommandFunc func(client *Client, args []string)

type CommandHandler struct {
	commands map[string]CommandFunc
}

func NewCommandHandler() *CommandHandler {
	ch := &CommandHandler{commands: make(map[string]CommandFunc)}

	RegisterLobbyCommands(ch)
	RegisterRoomCommands(ch)

	return ch
}

func (ch *CommandHandler) Register(cmd string, fn CommandFunc) {
	ch.commands[cmd] = fn
}

func (ch *CommandHandler) HandleCommand(c *Client, cmd string, args []string) bool {
	handler, ok := ch.commands[cmd]
	if !ok {
		var room string
		if c.state == "lobby" {
			room = "lobby"
		} else {
			room = c.currentRoom.name
		}

		out := NewOutgoing(
			"error",
			"server",
			room,
			fmt.Sprintf("Unknown command: %s", cmd),
		)
		c.writeCh <- out
		return true
	}

	handler(c, args)
	return true
}
