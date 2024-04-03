package main

import (
	"flag"
	"fmt"
	"log/slog"

	"github.com/GPlaczek/canvas/pkg/server"
)

func main() {
	bindAddr := flag.String("addr", "0.0.0.0", "Server listen address")
	listenPort := flag.String("port", "9090", "Server listen port")
	logLevel := flag.String("log-level", "INFO", "Log level")
	flag.Parse()

	var lv slog.Level
	switch *logLevel {
	case "DEBUG":
		lv = slog.LevelDebug
	case "INFO":
		lv = slog.LevelInfo
	case "WARN":
		lv = slog.LevelWarn
	case "ERROR":
		lv = slog.LevelError
	}

	cs := server.NewCanvasServer(lv)

	cs.StartServer(fmt.Sprintf("%s:%s", *bindAddr, *listenPort))
	fmt.Println("Hello world!")
}
