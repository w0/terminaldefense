package main

import (
	"crypto/sha1"
	"encoding/hex"
	"io"
	"log"
	"net/http"
	"os"
	"path"
)

func main() {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalln(err)
	}

	hub := &Hub{
		broadcast:  make(chan []byte),
		cmd:        make(chan []byte),
		clients:    make(map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}

	go hub.run()

	term := &Terminal{
		hub: hub,
	}

	go term.start()

	srvPath := path.Join(cwd, "static")

	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir(srvPath)))
	mux.HandleFunc("/ws", hub.handleWS)

	srv := http.Server{
		Addr:    "0.0.0.0:8080",
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
