package server

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/websocket"
)

type canvasServer struct {
	upgrader websocket.Upgrader
}

func (cs canvasServer) JoinCanvas(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s new connection from %s", r.Method, r.RemoteAddr)
	rid := r.PathValue("id")
	c, err := cs.upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Printf("could not upgrade connection (error %s)", err)
		return
	}

	cc := newClient(c)
	irid, err := strconv.Atoi(rid)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	cc.handleClient(irid)
}

func NewCanvasServer() canvasServer {
	cs := canvasServer{upgrader: websocket.Upgrader{}}

	mux := http.NewServeMux()
	mux.HandleFunc("/rooms/{id}", cs.JoinCanvas)

	http.Handle("/", mux)

	return cs
}

func (cs *canvasServer) StartServer(addr string) {
	log.Print("Starting server...")
	log.Fatal(http.ListenAndServe(addr, nil))
}
