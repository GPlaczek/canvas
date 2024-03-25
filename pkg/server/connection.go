package server

import (
	"fmt"

	"github.com/gorilla/websocket"
)

type canvasClient struct {
	connection *websocket.Conn
}

func newClient(conn *websocket.Conn) canvasClient {
	return canvasClient{connection: conn}
}

func (cc canvasClient) handleClient(rid int) {
	defer cc.connection.Close()

	cc.connection.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("Hello %d", rid)))
}
