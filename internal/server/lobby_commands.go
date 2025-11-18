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
	tmp := c.server.roomMgr.ListRooms()
	c.writeCh <- "Available Rooms:"
	for i, name := range tmp {
		c.writeCh <- fmt.Sprintf("%d. %s", i, name)
	}
}

func createRoomCommand(c *Client, args []string) {
	if len(args) == 1 {
		room, err := c.server.roomMgr.CreateRoom(args[0])
		if err != nil {
			c.writeCh <- err.Error()
			return
		}
		c.writeCh <- fmt.Sprintf("Room created with name %s", room.name)
		return
	}
	c.writeCh <- "Invalid argument. Space is not allowed in room names."
}

func joinRoomCommand(c *Client, args []string) {
	if len(args) == 1 {
		r := c.server.roomMgr.GetRoom(args[0])
		if r == nil {
			c.writeCh <- "Invalid room name."
			return
		}
		c.state = "inRoom"
		c.currentRoom = r
	}
	c.writeCh <- "Invalid use of command."
}

func lobbyHelpCommand(c *Client, args []string) {
	c.writeCh <- "Available commands: \n1. /rooms \n2. /create \n3. /join \n4. /exit"
}

func exitCommand(c *Client, args []string) {
	close(c.writeCh)
}
