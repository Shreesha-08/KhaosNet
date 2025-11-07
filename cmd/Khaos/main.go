package main

import (
	"KhaosNet/internal/server"
	"fmt"
	"net"
	"os"
)

func main() {
	listener, err := net.Listen("tcp", ":8182")
	if err != nil {
		fmt.Println("Failed to start the server: ", err.Error())
		os.Exit(1)
	}
	defer listener.Close()
	newServer := &server.Server{}
	newServer.NewServer(listener)
	newServer.AcceptConnection()
}
