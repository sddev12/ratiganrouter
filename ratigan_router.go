package ratiganrouter

import (
	"encoding/json"
	"io"
	"log"

	"golang.org/x/net/websocket"
)

type Event struct {
	Name string          `json:"name"`
	Data json.RawMessage `json:"data"`
}

type EventHandler func(ws *websocket.Conn, data json.RawMessage)

type Router struct {
	Routes map[string]EventHandler
}

func NewRouter() *Router {
	return &Router{
		make(map[string]EventHandler),
	}
}

func (r *Router) Handle(event string, handler EventHandler) {
	r.Routes[event] = handler
}

func (r *Router) Server(ws *websocket.Conn) {
	for {
		var event Event
		if err := websocket.JSON.Receive(ws, &event); err != nil {
			if err == io.EOF {
				log.Println("Connection Closed")
				return
			}
			log.Println("error receiving message: ", err)
			continue
		}

		if handler, ok := r.Routes[event.Name]; ok {
			handler(ws, event.Data)
		} else {
			log.Println("error no handler for event: ", event.Name)
		}
	}
}
