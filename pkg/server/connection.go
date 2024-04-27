package server

import (
	"log/slog"

	"google.golang.org/protobuf/proto"

	"github.com/gorilla/websocket"

	"github.com/GPlaczek/canvas/pkg/protocol"
)

type canvasClient struct {
	connection *websocket.Conn
	logger     *slog.Logger

	readMsg    chan *protocol.Message
	disconnect chan struct{}
	pause      chan struct{}
}

func NewClient(conn *websocket.Conn, logger *slog.Logger) *canvasClient {
	return &canvasClient{
		connection: conn,
		logger:     logger,
		readMsg:    make(chan *protocol.Message),
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
			cc.logger.Debug("Could not read message from the websocket", "error", err)
			break
		}

		var m protocol.Message
		err = proto.Unmarshal(msg, &m)
		if err != nil {
			cc.logger.Warn("Malformed message", "error", err)
			continue
		}

		cc.readMsg <- &m
	}
}

func (cc *canvasClient) HandleClient(room *Room) {
	defer func() {
		room.removeClient(cc)
		cc.connection.Close()
	}()

	go cc.readSocket()
	for {
		select {
		case mesg := <-cc.readMsg:
			switch mesg.Mtype {
			case protocol.MessageType_MESSAGE_POINT:
				if mesg.Line == nil || len(mesg.Line.Points) == 0 {
					continue
				}

				pt := ProtocolToPoint(mesg.Line.Points[0]) 

				err := room.addPoint(cc, *pt)
				if err != nil {
					cc.logger.Warn("Error processing addPoint message", "error", err)
				}
			case protocol.MessageType_MESSAGE_STOP:
				err := room.endLine(cc)
				if err != nil {
					cc.logger.Warn("Error processing endLine message", "error", err)
				}
			case protocol.MessageType_MESSAGE_CLEAN:
				cc.logger.Info("Cleaning the room")
				go room.cleanCanvas()
			default:
				cc.logger.Warn("Invalid message type")
			}
		case <-cc.pause:
			<-cc.pause
		case <-cc.disconnect:
			cc.logger.Info("Disconnecting client")
			return
		}
	}
}
