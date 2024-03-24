package server

import (
    "github.com/gorilla/websocket"
)

type canvasClient struct {
    connection *websocket.Conn
}

func newClient(conn *websocket.Conn) canvasClient {
    return canvasClient{connection: conn}
}

func (cc canvasClient) handleClient() {
    defer cc.connection.Close()

    cc.connection.WriteMessage(websocket.TextMessage, []byte("Hello"))
}
