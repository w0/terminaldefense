package main

import (
	"log"
	"net/http"
	"os"
	"path"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

type Hub struct {
	broadcast  chan []byte
	clients    map[*Client]bool
	register   chan *Client
	unregister chan *Client
}

type Client struct {
	addr string
	conn *websocket.Conn
	send chan []byte
	hub  *Hub
}

func main() {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalln(err)
	}

	hub := Hub{
		broadcast:  make(chan []byte),
		clients:    make(map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}

	go hub.run()

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
		conn: conn,
		send: make(chan []byte, 256),
		hub:  h,
	}

	h.register <- client

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

		log.Printf("%s, msg: %s", c.conn.RemoteAddr(), msg)
	}
}

func (c *Client) writePump() {
	for {
		msg := <-c.send
		c.conn.WriteMessage(websocket.TextMessage, msg)
	}
}
