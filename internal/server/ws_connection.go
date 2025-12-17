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

func NewWebSocketConn(connection *websocket.Conn) *WebSocketConn {
	return &WebSocketConn{conn: connection}
}

func (ws *WebSocketConn) Read() (string, error) {
	var msg IncomingMessage
	ctx, cancel := context.WithTimeout(context.Background(), time.Hour*1)
	defer cancel()

	if err := wsjson.Read(ctx, ws.conn, &msg); err != nil {
		return "", err
	}
	return msg.Text, nil
}

func (ws *WebSocketConn) ReadAndGetData() (*IncomingMessage, error) {
	var msg IncomingMessage
	ctx, cancel := context.WithTimeout(context.Background(), time.Hour*1)
	defer cancel()

	if err := wsjson.Read(ctx, ws.conn, &msg); err != nil {
		return nil, err
	}
	return &msg, nil
}

func (ws *WebSocketConn) Write(v interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return wsjson.Write(ctx, ws.conn, v)
}

func (ws *WebSocketConn) Close() error {
	return ws.conn.Close(websocket.StatusNormalClosure, "closing")
}
