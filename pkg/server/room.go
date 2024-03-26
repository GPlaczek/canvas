package server

import (
	"errors"
)

type Room struct {
	lines        [][]Point
	currentLines map[*canvasClient]*[]Point
}

func NewRoom() Room {
	return Room{
		lines:        make([][]Point, 0),
		currentLines: make(map[*canvasClient]*[]Point),
	}
}

func (r *Room) addClient(conn *canvasClient) error {
	_, ok := r.currentLines[conn]

	if ok {
		return errors.New("Client is already in the room")
	}

	return nil
}

func (r *Room) addPoint(conn *canvasClient, pt Point) error {
	line, ok := r.currentLines[conn]

	if !ok {
		return errors.New("Client is not a member of the room")
	}

	var ln []Point
	if line == nil {
		ln = make([]Point, 0)
		r.lines = append(r.lines, ln)
	} else {
		ln = *line
	}

	ln = append(ln, pt)

	return nil
}

func (r *Room) endLine(conn *canvasClient) error {
	line, ok := r.currentLines[conn]
	if !ok || line == nil {
		return errors.New("Client is not drawing a line")
	}

	r.currentLines[conn] = nil

	return nil
}
