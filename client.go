package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 4096
)

var (
	upgrader = websocket.Upgrader{}
	newline  = []byte{'\n'}
	space    = []byte{' '}
)

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	hub *Hub

	sub string

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte

	// Last touch time
	lastAccess time.Time

	// response is a response object to reply with
	response Response
}

// response is an object used to reply to the webhook with.
type Response struct {
	StatusCode int    `json:"statusCode,omitempty"`
	Body       string `json:"responseBody,omitempty"`
	Reflect    bool   `json:"reflect,omitempty"`
}

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) writePump() {
	c.lastAccess = time.Now()
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:

			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued chat messages to the current websocket message.
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// serveWs handles websocket requests from the peer.
func serveWs(hub *Hub, sub string, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := &Client{
		hub:        hub,
		sub:        sub,
		conn:       conn,
		send:       make(chan []byte, maxMessageSize),
		lastAccess: time.Now(),
		response:   Response{StatusCode: http.StatusOK, Reflect: true},
	}
	client.hub.register <- client

	hub.subs[sub] = client

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump()
}
