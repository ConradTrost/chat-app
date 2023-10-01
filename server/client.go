package main

import (
	"encoding/json"

	"github.com/gorilla/websocket"
)

type Client struct {
	clientId string
	// the socket connection for this client
	socket *websocket.Conn
	// channel to receive messages from other clients
	receive chan Message
	// room that client is currently in
	room *room
}

type Message struct {
	ClientId string `json:"clientId"`
	Message  string `json:"message"`
}

func (c *Client) read() {
	defer c.socket.Close()

	for {
		_, msg, err := c.socket.ReadMessage()
		if err != nil {
			return
		}

		newMsg := Message{
			ClientId: c.clientId,
			Message:  string(msg),
		}
		c.room.forward <- newMsg
	}
}

func (c *Client) write() {
	defer c.socket.Close()

	for msg := range c.receive {
		msgBytes, _ := json.Marshal(msg)

		err := c.socket.WriteMessage(websocket.TextMessage, msgBytes)
		if err != nil {
			return
		}
	}
}
