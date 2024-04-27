package server

import (
	"github.com/GPlaczek/canvas/pkg/protocol"
)

type RegisterRoom struct {
	Name string `json:"name"`
}

type Point struct {
	X int
	Y int
}

type Line struct {
	Ind    int
	Points []Point
}

func NewLine(i int) Line {
	return Line{
		Ind:    i,
		Points: make([]Point, 0),
	}
}

func ProtocolToPoint(pt *protocol.Point) *Point {
	return &Point {
		X: int(pt.X),
		Y: int(pt.Y),
	}
}

func (p *Point) Protocol() *protocol.Point {
	return &protocol.Point{
		X: int32(p.X),
		Y: int32(p.Y),
	}
}

func ProtocolToLine(ln *protocol.Line) *Line {
	line := &Line{
		Ind: int(ln.Ind),
		Points: make([]Point, len(ln.Points)),
	}

	for i, pt := range ln.Points {
		line.Points[i] = *ProtocolToPoint(pt) 
	}

	return line
}

func (l *Line) Protocol() *protocol.Line {
	p := make([]*protocol.Point, len(l.Points))

	for i, pt := range l.Points {
		p[i] = pt.Protocol()
	}

	return &protocol.Line {
		Ind: int32(l.Ind),
		Points: p,
	}
}
