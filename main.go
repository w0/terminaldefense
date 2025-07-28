package main

import (
	"crypto/sha1"
	"encoding/hex"
	"io"
	"log"
	"net/http"
	"os"
	"path"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	port := os.Getenv("TD_PORT")
	log.Printf("srv on port %s", port)

	ex, err := os.Getwd()
	if err != nil {
		log.Fatalln(err)
	}

	staticDir := path.Join(ex, "static")

	log.Printf("srv static from %s", staticDir)

	hub := &Hub{
		action:     make(chan Action),
		broadcast:  make(chan []byte),
		cmd:        make(chan []byte),
		clients:    make(map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		nextRole:   "hacker",
	}

	go hub.run()

	term := &Terminal{
		hub:     hub,
		pending: make(map[string]*PendingCommand),
	}

	go term.start()

	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir(staticDir)))
	mux.HandleFunc("/ws", hub.handleWS)

	srv := http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	srv.ListenAndServe()
}

func hashIt(ip string) string {
	h := sha1.New()

	io.WriteString(h, ip)

	hStr := hex.EncodeToString(h.Sum(nil))

	strLen := len(hStr)

	return hStr[strLen-7:]
}
