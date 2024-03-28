package server

import (
	"encoding/json"
)

const (
	MESSAGE_POINT = 0
	MESSAGE_STOP  = 1
)

type Message struct {
	MType int   `json:"type"`
	Point Point `json:"point,omitempty"`
}

type Point struct {
	X int `json:"x"`
	Y int `json:"y"`
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

func (l Line) MarshalJSON() ([]byte, error) {
	var pts struct {
		Ind int   `json:"ind"`
		X   []int `json:"x"`
		Y   []int `json:"y"`
	}

	pts.Ind = l.Ind
	pts.X = make([]int, len(l.Points))
	pts.Y = make([]int, len(l.Points))
	for i, pt := range l.Points {
		pts.X[i] = pt.X
		pts.Y[i] = pt.Y
	}

	return json.Marshal(&pts)
}
