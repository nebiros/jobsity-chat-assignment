package websocket

import (
	"context"
	"log/slog"
	"time"

	"github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 15
	maxMessageSize = 1024 * 1024
)

var (
	Upgrader = websocket.Upgrader{
		ReadBufferSize:  4096,
		WriteBufferSize: 4096,
	}
)

type Client[T any] struct {
	ctx  context.Context
	hub  *Hub[T]
	conn *websocket.Conn
	send chan T
}

func NewClient[T any](ctx context.Context, hub *Hub[T], conn *websocket.Conn) *Client[T] {
	return &Client[T]{
		ctx:  ctx,
		hub:  hub,
		conn: conn,
		send: make(chan T, maxMessageSize),
	}
}

func (c *Client[T]) ReadPump(cb func(chan T, T)) {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	_ = c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		if err := c.conn.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
			slog.Error("pong error", slog.Any("error", err))

			return err
		}

		slog.Debug("pong")

		return nil
	})

	for {
		var message T

		err := c.conn.ReadJSON(&message)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				slog.Error("unexpected close error", slog.Any("message", message), slog.Any("error", err))
			} else {
				slog.Error("unable to read JSON", slog.Any("message", message), slog.Any("error", err))
			}

			break
		}

		go cb(c.hub.broadcast, message)

		c.hub.broadcast <- message
	}
}

func (c *Client[T]) WritePump() {
	ticker := time.NewTicker(pingPeriod)

	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			_ = c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				slog.Error("unable to get message from send channel", slog.Any("message", message))
				_ = c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := c.conn.WriteJSON(message); err != nil {
				slog.Error("unable to write JSON message", slog.Any("message", message), slog.Any("error", err))
				return
			}

		case <-ticker.C:
			_ = c.conn.SetWriteDeadline(time.Now().Add(writeWait))

			slog.Debug("ping")

			if err := c.conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				slog.Error("ping error", slog.Any("error", err))

				return
			}
		}
	}
}
