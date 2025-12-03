package server

import (
	"context"
	"time"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
)

type WebSocketConn struct {
	conn *websocket.Conn
}

type IncomingMessage struct {
	Text string `json:"text"`
}

func NewWebSocketConn(connection *websocket.Conn) *WebSocketConn {
	return &WebSocketConn{conn: connection}
}

func (ws *WebSocketConn) Read() (string, error) {
	var msg IncomingMessage
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	if err := wsjson.Read(ctx, ws.conn, &msg); err != nil {
		return "", err
	}
	return msg.Text, nil
}

func (ws *WebSocketConn) Write(msg string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return wsjson.Write(ctx, ws.conn, msg)
}

func (ws *WebSocketConn) Close() error {
	return ws.conn.Close(websocket.StatusNormalClosure, "closing")
}
