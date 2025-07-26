package main

import (
	"log"
	"net/http"
	"os"
	"path"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

func main() {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalln(err)
	}

	srvPath := path.Join(cwd, "static")

	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir(srvPath)))
	mux.HandleFunc("/ws", serverWs)

	srv := http.Server{
		Addr:    "0.0.0.0:8080",
		Handler: mux,
	}

	srv.ListenAndServe()
}

func serverWs(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}

	go doThings(ws)

}

func doThings(ws *websocket.Conn) {
	log.Printf("ws open from %s", ws.NetConn().RemoteAddr().String())

	defer ws.Close()
	ws.WriteMessage(websocket.TextMessage, []byte("hi from inside of a cool go routine!"))

}
