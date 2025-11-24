package main

import (
  "github.com/gorilla/websocket"
  "log"
  "encoding/json"
)

type Client struct {
  conn *websocket.Conn
  hub *Hub
  send chan *OutboundMessage
  handshake chan *Handshake 
  name string
  id string
}

type Handshake struct {
  ClientId string `json:"clientId"`
}

type InboundMessage struct {
  ClientId string `json:"clientId"`
  Data Coordinates `json:"data,omitempty"` 
  Type string `json:"type"`
  ServerId string
}

type Coordinates struct {
  X int `json:"x"`
  Y int `json:"y"`
}

type OutboundMessage struct {
  Data Coordinates `json:"data"` // should rename to payload
  Type string `json:"type"`
  ClientName string `json:"clientName"`
}

/*
Read stream of messages from clients and publish to redis channel
*/
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

    // unmarshal inbound message
    var msg InboundMessage
    if err := json.Unmarshal(message, &msg); err != nil {
      log.Printf("invalid message: %v", err)
      continue
    }

    // assign server ID
    msg.ServerId = c.hub.serverId

    // marshal it back to JSON
    jsonBytes, err := json.Marshal(msg)
    if err != nil {
      log.Printf("marshal error: %v", err)
      continue
    }

    // publish to Redis
    c.hub.redis.Publish(c.hub.ctx, "mouse_events", jsonBytes)

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
        continue
      }

      jsonBytes, _ := json.Marshal(message)
      err := c.conn.WriteMessage(websocket.TextMessage, jsonBytes)

      if err != nil {
        log.Println("WriteMessage error:", err)
        continue
      }
    case handshake := <-c.handshake:
      if (handshake == nil) {
        continue
      }
      w,_ := c.conn.NextWriter(websocket.TextMessage)
      payload := map[string]interface{}{
        "type":     "handshake",
        "clientId": handshake.ClientId,
      }
      b, _ := json.Marshal(payload)
      w.Write(b)
      w.Close()
    }
  }
}
