package server

import (
	"encoding/json"
	"log/slog"

	"github.com/gorilla/websocket"
)

type canvasClient struct {
	connection *websocket.Conn
	logger     *slog.Logger

	readMsg    chan Message
	disconnect chan struct{}
	pause      chan struct{}
}

func NewClient(conn *websocket.Conn, logger *slog.Logger) *canvasClient {
	return &canvasClient{
		connection: conn,
		logger:     logger,
		readMsg:    make(chan Message),
		disconnect: make(chan struct{}),
		pause:      make(chan struct{}),
	}
}

func (cc *canvasClient) readSocket() {
	defer func() {
		cc.disconnect <- struct{}{}
	}()

	for {
		_, msg, err := cc.connection.ReadMessage()
		if err != nil {
			break
		}

		var m Message
		err = json.Unmarshal(msg, &m)
		if err != nil {
			cc.logger.Warn("Malformed message", "error", err)
			continue
		}

		cc.readMsg <- m
	}
}

func (cc *canvasClient) HandleClient(room *Room) {
	defer cc.connection.Close()

	go cc.readSocket()
	for {
		select {
		case mesg := <-cc.readMsg:
			switch mesg.MType {
			case MESSAGE_POINT:
				if len(mesg.Line.Points) == 0 {
					continue
				}

				err := room.addPoint(cc, mesg.Line.Points[0])
				if err != nil {
					cc.logger.Warn("Error processing addPoint message", "error", err)
				}
			case MESSAGE_STOP:
				err := room.endLine(cc)
				if err != nil {
					cc.logger.Warn("Error processing endLine message", "error", err)
				}
			case MESSAGE_CLEAN:
				cc.logger.Info("Cleaning the room")
				go room.cleanCanvas()
			default:
				cc.logger.Warn("Invalid message type")
			}
		case <-cc.pause:
			<-cc.pause
		case <-cc.disconnect:
			break
		}
	}
}
