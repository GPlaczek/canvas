package server

import (
	"errors"
	"log"
)

type Room struct {
	lines        []Line
	currentLines map[*canvasClient]int
}

func NewRoom() Room {
	return Room{
		lines:        make([]Line, 0),
		currentLines: make(map[*canvasClient]int),
	}
}

func (r *Room) addClient(conn *canvasClient) error {
	_, ok := r.currentLines[conn]

	if ok {
		return errors.New("Client is already in the room")
	}

	i := 0
	for _, line := range r.lines {
		err := conn.connection.WriteJSON(line)
		if err != nil {
			log.Println("Could not serialize the line")
		}
		i++
	}

	log.Println("Adding new client to a room")
	r.currentLines[conn] = -1

	return nil
}

func (r *Room) addPoint(conn *canvasClient, pt Point) error {
	line, ok := r.currentLines[conn]

	if !ok {
		return errors.New("Client is not a member of the room")
	}

	var ln *Line
	if line == -1 {
		r.currentLines[conn] = len(r.lines)
		r.lines = append(r.lines, NewLine(len(r.lines)))
		ln = &r.lines[len(r.lines)-1]
	} else {
		ln = &r.lines[line]
	}

	ln.Points = append(ln.Points, pt)

	for client := range r.currentLines {
		if client != conn {
			client.connection.WriteJSON(Line{
				Ind:    line,
				Points: []Point{pt},
			})
		}
	}
	return nil
}

func (r *Room) endLine(conn *canvasClient) error {
	line, ok := r.currentLines[conn]
	if !ok || line == -1 {
		return errors.New("Client is not drawing a line")
	}

	r.currentLines[conn] = -1

	return nil
}
