package main

import (
	"github.com/gorilla/websocket"
)

type Client struct {
	name string
	conn *websocket.Conn
	send chan []byte
	hub  *Hub
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
