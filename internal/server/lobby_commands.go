package server

import "fmt"

func RegisterLobbyCommands(ch *CommandHandler) {
	ch.Register("/rooms", roomsCommand)
	ch.Register("/create", createRoomCommand)
	ch.Register("/join", joinRoomCommand)
	ch.Register("/help", lobbyHelpCommand)
	ch.Register("/exit", exitCommand)
}

func roomsCommand(c *Client, args []string) {
	rooms := c.server.roomMgr.ListRooms()
	resp := NewOutgoing(
		"rooms_list",
		"server",
		"lobby",
		"",
	)
	resp.Data = map[string]any{
		"rooms": rooms,
	}
	c.writeCh <- resp
}

func createRoomCommand(c *Client, args []string) {
	if len(args) != 1 {
		out := NewOutgoing(
			"error",
			"server",
			"lobby",
			"Invalid argument. Usage: /create <roomname>",
		)
		c.writeCh <- out
		return
	}

	roomName := args[0]
	room, err := c.server.roomMgr.CreateRoom(roomName)
	if err != nil {
		out := NewOutgoing(
			"error",
			"server",
			"lobby",
			err.Error(),
		)
		c.writeCh <- out
		return
	}

	out := NewOutgoing(
		"room_created",
		"server",
		room.name,
		fmt.Sprintf("Room created: %s", room.name),
	)
	c.writeCh <- out
}

func joinRoomCommand(c *Client, args []string) {
	if len(args) != 1 {
		out := NewOutgoing(
			"error",
			"server",
			"lobby",
			"Usage: /join <roomname>",
		)
		c.writeCh <- out
		return
	}

	roomName := args[0]
	room := c.server.roomMgr.GetRoom(roomName)
	if room == nil {
		out := NewOutgoing(
			"error",
			"server",
			"lobby",
			"Room not found.",
		)
		c.writeCh <- out
		return
	}

	c.state = "inRoom"
	c.currentRoom = room

	ack := NewOutgoing(
		"system",
		"server",
		roomName,
		fmt.Sprintf("Joined room: %s", roomName),
	)
	c.writeCh <- ack
	c.currentRoom.broadcaster.joinCh <- c
}

func lobbyHelpCommand(c *Client, args []string) {
	helpText := "Available commands:\n" +
		"1. /rooms\n" +
		"2. /create <room>\n" +
		"3. /join <room>\n" +
		"4. /name <newname>\n" +
		"5. /exit"

	out := NewOutgoing(
		"system",
		"server",
		"lobby",
		helpText,
	)
	c.writeCh <- out
}

func exitCommand(c *Client, args []string) {
	close(c.writeCh)
}
