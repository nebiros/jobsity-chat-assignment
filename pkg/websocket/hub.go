package websocket

import (
	"context"
	"log/slog"
)

type Hub[T any] struct {
	ctx context.Context

	// clients holds registered clients.
	clients map[*Client[T]]bool

	// broadcast inbound messages from the clients.
	broadcast chan T

	// register requests from the clients.
	register chan *Client[T]

	// unregister requests from clients.
	unregister chan *Client[T]
}

func NewHub[T any](ctx context.Context) *Hub[T] {
	return &Hub[T]{
		ctx:        ctx,
		clients:    make(map[*Client[T]]bool),
		broadcast:  make(chan T),
		register:   make(chan *Client[T]),
		unregister: make(chan *Client[T]),
	}
}

func (h *Hub[T]) Register() chan *Client[T] {
	return h.register
}

func (h *Hub[T]) Unregister() chan *Client[T] {
	return h.unregister
}

func (h *Hub[T]) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true

		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				slog.Debug("de registering client", slog.Any("client", client))

				delete(h.clients, client)
				close(client.send)
			}

		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.send <- message:
					slog.Debug("broadcasting message", slog.Any("message", message))

				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}
