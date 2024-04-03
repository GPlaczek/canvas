package server

import (
	"errors"
	"log/slog"
	"sync"
)

type Room struct {
	lines        []Line
	currentLines sync.Map /* (*canvasClient, int) */
	logger       *slog.Logger

	linesMtx sync.RWMutex
}

func NewRoom(logger *slog.Logger) Room {
	return Room{
		logger:       logger,
		lines:        make([]Line, 0),
		currentLines: sync.Map{},
	}
}

func (r *Room) addClient(conn *canvasClient) error {
	_, loaded := r.currentLines.LoadOrStore(conn, -1)

	if loaded {
		return errors.New("Client is already in the room")
	}

	r.linesMtx.RLock()
	defer r.linesMtx.RUnlock()

	i := 0
	for _, line := range r.lines {
		err := conn.connection.WriteJSON(Message{
			MType: MESSAGE_POINT,
			Line:  line,
		})
		if err != nil {
			r.logger.Warn("Could not serialize the line")
		}
		i++
	}

	return nil
}

func (r *Room) removeClient(conn *canvasClient) {
	r.currentLines.Delete(conn)
}

func (r *Room) addPoint(conn *canvasClient, pt Point) error {
	__line, ok := r.currentLines.Load(conn)
	line := __line.(int)

	if !ok {
		return errors.New("Client is not a member of the room")
	}

	var ln *Line
	if line == -1 {
		r.linesMtx.Lock()
		line = len(r.lines)
		r.lines = append(r.lines, NewLine(line))
		r.linesMtx.Unlock()

		r.currentLines.Store(conn, line)
		ln = &r.lines[line]
	} else {
		ln = &r.lines[line]
	}

	ln.Points = append(ln.Points, pt)

	r.currentLines.Range(func(client, _ any) bool {
		if client != conn {
			err := client.(*canvasClient).connection.WriteJSON(Message{
				Line: Line{
					Ind:    line,
					Points: []Point{pt},
				},
				MType: MESSAGE_POINT,
			})

			if err != nil {
				r.logger.Warn("Could not sent lines to a client", "error", err)
			}
		}
		return true
	})

	return nil
}

func (r *Room) endLine(conn *canvasClient) error {
	line, ok := r.currentLines.Load(conn)
	if !ok || line == -1 {
		return errors.New("Client is not drawing a line")
	}

	r.currentLines.Store(conn, -1)

	return nil
}

func (r *Room) cleanCanvas() {
	pauseUnpause := func() {
		r.currentLines.Range(func(client, _ any) bool {
			// TODO: run a goroutine for each client so as not to
			// block waiting for each one to receive the message
			client.(*canvasClient).pause <- struct{}{}
			return true
		})
	}

	pauseUnpause()
	defer pauseUnpause()

	r.lines = make([]Line, 0)
	r.currentLines.Range(func(client, _ any) bool {
		r.currentLines.Store(client, -1)
		err := client.(*canvasClient).connection.WriteJSON(Message{
			MType: MESSAGE_CLEAN,
		})

		if err != nil {
			r.logger.Warn("Could not send clean message to a client", "error", err)
		}
		return true
	})
}
