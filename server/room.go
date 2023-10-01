package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type room struct {
	clients map[*Client]bool
	join    chan *Client
	leave   chan *Client
	forward chan Message
}

func newRoom() *room {
	return &room{
		forward: make(chan Message),
		join:    make(chan *Client),
		leave:   make(chan *Client),
		clients: make(map[*Client]bool),
	}
}

func (r *room) run() {
	for {
		select {
		case client := <-r.join:
			r.clients[client] = true
		case client := <-r.leave:
			log.Println("Client is leaving:", client.clientId)
			delete(r.clients, client)
			close(client.receive)
		case msg := <-r.forward:
			for client := range r.clients {
				client.receive <- msg
			}
		}
	}
}

var upgrader = &websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	}}

func (r *room) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	params := req.URL.Query()
	username := params.Get("username")

	if username == "" {
		return
	}

	socket, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Fatal(err)
		// return
	}

	client := &Client{
		clientId: username,
		socket:   socket,
		receive:  make(chan Message, 256),
		room:     r,
	}
	r.join <- client

	defer func() { r.leave <- client }()

	go client.write()
	client.read()
}
