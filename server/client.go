package main

import (
    "github.com/gorilla/websocket"
    "log"
)

type Client struct {
  conn *websocket.Conn
  hub *Hub
  send chan *Message
  name string
}

type Message struct {
  client *Client
  message []byte
}

func (c *Client) readPump() {
  defer func() {
    c.hub.unregister <- c
    c.conn.Close()
  }()

  for {
    _, message, err := c.conn.ReadMessage()
    if err != nil {
      log.Printf("error: %v", err)
      return
    }
    messageObj := &Message {
      client: c,
      message: message,
    }
    c.hub.broadcast <- messageObj 
  }
}

func (c *Client) writePump() {
  defer func() {
    c.hub.unregister <- c
    c.conn.Close()
  }()

  for {
    select {
    case message, ok := <-c.send:
      if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

      w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message.message)

			if err := w.Close(); err != nil {
				return
			}
    }
  }
}
