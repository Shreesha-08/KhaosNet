package server

import "net"

type Client struct {
	conn    net.Conn
	name    string
	writeCh chan string
	doneCh  chan struct{}
}
