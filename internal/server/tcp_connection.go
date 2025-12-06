package server

import (
	"net"
	"strings"
)

type TCPConn struct {
	conn net.Conn
}

func NewTCPConnection(connection net.Conn) *TCPConn {
	return &TCPConn{conn: connection}
}

func (tc *TCPConn) Read() (string, error) {
	buf := make([]byte, 256)
	n, err := tc.conn.Read(buf)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(buf[:n])), nil
}

func (tc *TCPConn) ReadAndGetData() (*IncomingMessage, error) {
	// To be implemented
	return nil, nil
}

func (tc *TCPConn) Write(msg string) error {
	_, err := tc.conn.Write([]byte(msg + "\n"))
	return err
}

func (tc *TCPConn) Close() error {
	return tc.conn.Close()
}
