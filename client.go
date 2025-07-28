package main

import (
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
)

type Client struct {
	name string
	conn *websocket.Conn
	send chan []byte
	hub  *Hub
	role string // hacker, sysadmin
}

type Message struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}

type Command struct {
	Command string `json:"command"`
}

type Action struct {
	Action string `json:"action"`
	Id     string `json:"id"`
}

func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	var msg Message

	for {
		err := c.conn.ReadJSON(&msg)
		if err != nil {
			log.Printf("json read %s", err)
		}

		switch msg.Type {
		case "CMD":
			var cmdMsg Command
			err := json.Unmarshal(msg.Data, &cmdMsg)
			if err != nil {
				log.Printf("cmd %s", err)
			}
			c.hub.cmd <- []byte(cmdMsg.Command)
		case "ACTION":
			var actionMsg Action
			err := json.Unmarshal(msg.Data, &actionMsg)
			if err != nil {
				log.Printf("action %s", err)
			}
			c.hub.action <- actionMsg
		}
	}
}

func (c *Client) writePump() {
	for {
		msg := <-c.send
		c.conn.WriteMessage(websocket.TextMessage, msg)
	}
}
