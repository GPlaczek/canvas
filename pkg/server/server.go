package server

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type canvasServer struct {
	upgrader  websocket.Upgrader
	roomsLock sync.RWMutex
	rooms     map[string]*Room
}

func (cs *canvasServer) JoinCanvas(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s new connection from %s", r.Method, r.RemoteAddr)

	rid := r.PathValue("id")

	switch r.Method {
	case http.MethodGet:
		// TODO: this lock does not need to be up for the connection upgrading
		cs.roomsLock.RLock()
		defer cs.roomsLock.RUnlock()

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
		room.addClient(&cc)
		go cc.HandleClient(room)
	case http.MethodPost:
		cs.roomsLock.Lock()
		defer cs.roomsLock.Unlock()

		_, ok := cs.rooms[rid]
		if ok {
			w.WriteHeader(http.StatusConflict)
			return
		}

		nr := NewRoom()
		cs.rooms[rid] = &nr
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}

func NewCanvasServer() *canvasServer {
	cs := canvasServer{
		upgrader: websocket.Upgrader{
			CheckOrigin: func(_ *http.Request) bool { return true },
		},
		rooms: make(map[string]*Room),
	}

	mux := http.NewServeMux()
	// mux.HandleFunc("/room", cs.RoomsHandler)
	mux.HandleFunc("/room/{id}", cs.JoinCanvas)

	http.Handle("/", mux)

	return &cs
}

func (cs *canvasServer) StartServer(addr string) {
	log.Print("Starting server...")
	log.Fatal(http.ListenAndServe(addr, nil))
}
