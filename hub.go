package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

type Hub struct {
	action     chan Action
	broadcast  chan []byte
	cmd        chan []byte
	previous   []byte
	clients    map[*Client]bool
	register   chan *Client
	unregister chan *Client
	terminal   *Terminal
	nextRole   string // this should be a map of active roles and count.
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
		role: h.nextRole,
	}

	// need logic to balanace roles across pool of players
	// hub should balance based on active players, move clients between roles to keep game moving.
	if h.nextRole == "hacker" {
		h.nextRole = "sysadmin"
	} else {
		h.nextRole = "hacker"
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
