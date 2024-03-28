package server

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
	Ind    int     `json:"ind"`
	Points []Point `json:"point"`
}

func NewLine(i int) Line {
	return Line{
		Ind:    i,
		Points: make([]Point, 0),
	}
}
