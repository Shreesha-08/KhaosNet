package server

import "strings"

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

func (ch *CommandHandler) HandleCommand(c *Client, input string) bool {
	parts := strings.Fields(input)
	if len(parts) == 0 {
		return false
	}

	cmd := parts[0]
	args := parts[1:]

	handler, ok := ch.commands[cmd]
	if !ok {
		c.writeCh <- "Unknown command"
		return true
	}

	handler(c, args)
	return true
}
