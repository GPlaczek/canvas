package server

import (
    "log"
    "net/http"

    "github.com/gorilla/websocket"
)

type canvasServer struct {
    upgrader websocket.Upgrader
}

func (cs canvasServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    log.Printf("new connection from %s", r.RemoteAddr)
    c, err := cs.upgrader.Upgrade(w, r, nil)

    if err != nil {
        log.Printf("could not upgrade connection (error %s)", err)
        return
    }

    cc := newClient(c)
    cc.handleClient()
}

func NewCanvasServer() canvasServer {
    cs := canvasServer{upgrader: websocket.Upgrader{}}
    http.Handle("/", cs)
    return cs
}

func (cs *canvasServer) StartServer(addr string) {
    log.Print("Starting server...")
    log.Fatal(http.ListenAndServe(addr, nil))
}
