package main

import (
	"log"
	"math/rand/v2"
	"net/http"
	"os"
	"os/exec"
	"path"
	"time"

	"github.com/creack/pty"
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

	for {
		size := rand.IntN(64)
		buf := []byte("This is a really long message that I have written. You wont see it all!")

		err := ws.WriteMessage(websocket.TextMessage, buf[:size])

		if err != nil {
			break
		}

		time.Sleep(time.Second * 2)

		mt, p, err := ws.ReadMessage()

		log.Println(mt)
		log.Println(string(p))
		if err != nil {
			log.Printf("error reading ws msg %s", err)
		}

	}

}

func startShell() ([]byte, error) {
	ptmx, err := pty.Start(exec.Command("bash", "-c", "ls"))
	if err != nil {
		log.Fatal(err)
	}

	buf := make([]byte, 1024)
	r, err := ptmx.Read(buf)
	if err != nil {
		log.Println(err)
	}

	return buf[:r], nil

}

func doTick() {

}
