package server

import (
	"errors"
	"log"
	"sync"
)

type Room struct {
	lines        []Line
	currentLines sync.Map  /* (*canvasClient, int) */

	linesMtx     sync.RWMutex
}

func NewRoom() Room {
	return Room{
		lines:        make([]Line, 0),
		currentLines: sync.Map{},
	}
}

func (r *Room) addClient(conn *canvasClient) error {
	_, loaded := r.currentLines.LoadOrStore(conn, -1)

	if !loaded {
		return errors.New("Client is already in the room")
	}

	log.Println("Adding new client to a room")

	r.linesMtx.RLock()
	defer r.linesMtx.Unlock()

	i := 0
	for _, line := range r.lines {
		err := conn.connection.WriteJSON(line)
		if err != nil {
			log.Println("Could not serialize the line")
		}
		i++
	}

	return nil
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

	r.currentLines.Range(func (client, _ any) bool {
		if client != conn {
			client.(*canvasClient).connection.WriteJSON(Line{
				Ind:    line,
				Points: []Point{pt},
			})
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
