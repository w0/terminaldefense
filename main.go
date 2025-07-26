package main

import (
	"log"
	"math/rand/v2"
	"net/http"
	"os"
	"path"
	"time"

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

	go writePump(ws)
	go readPump(ws)

}

func writePump(ws *websocket.Conn) {
	log.Printf("ws open: %s", ws.NetConn().LocalAddr().String())

	defer ws.Close()

	buf := []byte("This is a really long message that I have written. You wont see it all!")

	for {
		size := rand.IntN(len(buf))

		err := ws.WriteMessage(websocket.TextMessage, buf[:size])
		if err != nil {
			log.Printf("ws write error: %s", err)
		}

		time.Sleep(time.Millisecond * 60)

	}

}

func readPump(ws *websocket.Conn) {
	for {
		_, p, err := ws.ReadMessage()
		if err != nil {
			log.Printf("ws read err: %s", err)
		}

		log.Printf("ws msg: %s", string(p))
	}
}
