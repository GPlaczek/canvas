package main

import (
	"fmt"

	"github.com/GPlaczek/canvas/pkg/server"
)

func main() {
	cs := server.NewCanvasServer()
	cs.StartServer("localhost:9090")
	fmt.Println("Hello world!")
}
