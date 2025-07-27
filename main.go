package main

import (
	"crypto/sha1"
	"encoding/hex"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"

	"github.com/creack/pty"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

type Hub struct {
	broadcast  chan []byte
	cmd        chan []byte
	previous   []byte
	clients    map[*Client]bool
	register   chan *Client
	unregister chan *Client
}

type Client struct {
	name string
	conn *websocket.Conn
	send chan []byte
	hub  *Hub
}

type Terminal struct {
	pty *os.File
	hub *Hub
}

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

func (h *Hub) handleWS(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("conn upgrade failed %s", err)
	}

	client := &Client{
		name: hashIt(conn.NetConn().RemoteAddr().String()),
		conn: conn,
		send: make(chan []byte, 256),
		hub:  h,
	}

	h.register <- client

	client.send <- h.previous

	go client.readPump()
	go client.writePump()

}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
			log.Printf("client %s connected, total: %d", client.conn.RemoteAddr().String(), len(h.clients))

		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
				log.Printf("client %s disconnected", client.conn.RemoteAddr().String())
			}

		case message := <-h.broadcast:
			h.previous = message
			for client := range h.clients {
				client.send <- message
			}
		}
	}
}

func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	for {
		_, msg, err := c.conn.ReadMessage()
		if err != nil {
			break
		}

		c.hub.cmd <- msg
	}
}

func (c *Client) writePump() {
	for {
		msg := <-c.send
		c.conn.WriteMessage(websocket.TextMessage, msg)
	}
}

func hashIt(ip string) string {
	h := sha1.New()

	io.WriteString(h, ip)

	hStr := hex.EncodeToString(h.Sum(nil))

	strLen := len(hStr)

	return hStr[strLen-7:]
}

func (t *Terminal) start() {
	pty, err := pty.Start(exec.Command("bash"))
	if err != nil {
		log.Fatalf("bash start %s", err)
	}

	t.pty = pty

	go t.readFromShell()
	go t.writeToShell()

}

func (t *Terminal) readFromShell() {
	buf := make([]byte, 1024)

	for {
		n, err := t.pty.Read(buf)
		if err != nil {
			log.Printf("terminal read failed %s", err)
		}

		t.hub.broadcast <- buf[:n]
	}

}

func (t *Terminal) writeToShell() {
	for {
		msg := <-t.hub.cmd

		msg = append(msg, '\n')

		_, err := t.pty.Write(msg)
		if err != nil {
			log.Panicln(err)
		}

	}
}
