package server

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"sync"

	"github.com/gorilla/websocket"
)

type canvasServer struct {
	upgrader  websocket.Upgrader
	roomsLock sync.RWMutex
	rooms     map[string]*Room
	logger    *slog.Logger
}

func (cs *canvasServer) JoinCanvas(w http.ResponseWriter, r *http.Request) {
	addrLogger := cs.logger.With("address", r.RemoteAddr)
	addrLogger.Info("New connection", "method", r.Method)

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
			addrLogger.Warn("could not upgrade connection", "error", err)
			return
		}

		cc := NewClient(c, addrLogger)
		err = room.addClient(cc)
		if err != nil {
			addrLogger.Warn("Client already in the room")
			w.WriteHeader(http.StatusConflict)
			return
		}
		go cc.HandleClient(room)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}

func (cs *canvasServer) RoomsHandler(w http.ResponseWriter, r *http.Request) {
	addrLogger := cs.logger.With("address", r.RemoteAddr)
	addrLogger.Info("New connection", "method", r.Method)

	switch r.Method {
	case http.MethodPost:
		var rr RegisterRoom
		err := json.NewDecoder(r.Body).Decode(&rr)
		if err != nil {
			addrLogger.Info("Bad register room request", "error", err)
			w.WriteHeader(http.StatusUnprocessableEntity)
			return
		}

		if rr.Name == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		cs.roomsLock.Lock()
		defer cs.roomsLock.Unlock()

		_, ok := cs.rooms[rr.Name]
		if ok {
			w.WriteHeader(http.StatusConflict)
			return
		}

		nr := NewRoom(cs.logger)
		cs.rooms[rr.Name] = &nr
	case http.MethodGet:
		cs.roomsLock.RLock()
		defer cs.roomsLock.RUnlock()

		rl := make([]string, len(cs.rooms))

		i := 0
		for k := range cs.rooms {
			rl[i] = k
			i++
		}

		body, err := json.Marshal(rl)
		if err != nil {
			addrLogger.Error("Could not serialize response body")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err = w.Write(body)
		if err != nil {
			addrLogger.Error("Could not write response body")
		}
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}

func NewCanvasServer(level slog.Level) *canvasServer {
	lv := new(slog.LevelVar)
	lv.Set(level)

	cs := canvasServer{
		upgrader: websocket.Upgrader{
			CheckOrigin: func(_ *http.Request) bool { return true },
		},
		rooms: make(map[string]*Room),
		logger: slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: lv,
		})),
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/rooms", cs.RoomsHandler)
	mux.HandleFunc("/rooms/{id}", cs.JoinCanvas)

	http.Handle("/", mux)

	return &cs
}

func (cs *canvasServer) StartServer(addr string) {
	cs.logger.Info("Starting server...")
	cs.logger.Error("Server stopped", "error", http.ListenAndServe(addr, nil))
}
