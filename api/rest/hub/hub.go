package hub

import (
	"github.com/gorilla/websocket"
)

type Hub struct {
	Clients    map[string]map[*websocket.Conn]chan []byte
	Register   chan *Connect
	Unregister chan *Connect
}

type Connect struct {
	Conn *websocket.Conn
	Send chan []byte
	User string
}

// NewHub returns new Hub.
func NewHub() *Hub {
	return &Hub{
		Clients:    make(map[string]map[*websocket.Conn]chan []byte),
		Register:   make(chan *Connect),
		Unregister: make(chan *Connect),
	}
}

func (h Hub) Close() {
	close(h.Register)
	close(h.Unregister)
	for _, userClients := range h.Clients {
		for _, send := range userClients {
			close(send)
		}
	}
}
