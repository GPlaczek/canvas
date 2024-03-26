package server

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type canvasServer struct {
	upgrader websocket.Upgrader
	rooms    map[string]Room
}

func (cs canvasServer) JoinCanvas(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s new connection from %s", r.Method, r.RemoteAddr)

	rid := r.PathValue("id")

	switch r.Method {
	case http.MethodGet:
		{
			room, ok := cs.rooms[rid]
			if !ok {
				w.WriteHeader(http.StatusNotFound)
				return
			}

			c, err := cs.upgrader.Upgrade(w, r, nil)
			if err != nil {
				log.Printf("could not upgrade connection (error %s)", err)
				return
			}

			cc := NewClient(c)
			cc.HandleClient(&room)
		}
	case http.MethodPost:
		{
			_, ok := cs.rooms[rid]
			if ok {
				w.WriteHeader(http.StatusConflict)
				return
			}

			cs.rooms[rid] = NewRoom()
		}
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}

func NewCanvasServer() canvasServer {
	cs := canvasServer{
		upgrader: websocket.Upgrader{},
		rooms:    make(map[string]Room),
	}

	mux := http.NewServeMux()
	// mux.HandleFunc("/room", cs.RoomsHandler)
	mux.HandleFunc("/room/{id}", cs.JoinCanvas)

	http.Handle("/", mux)

	return cs
}

func (cs *canvasServer) StartServer(addr string) {
	log.Print("Starting server...")
	log.Fatal(http.ListenAndServe(addr, nil))
}
