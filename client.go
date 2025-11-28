package main

import (
  "log"
  "math/rand"
  "fmt"
  "net/http"

  "github.com/gorilla/websocket"
  "github.com/google/uuid"
	"google.golang.org/protobuf/proto"
	 events "github.com/walshyb/whiteboard/proto"
)

type Client struct {
  conn *websocket.Conn
  hub *Hub
  send chan *events.ServerMessage
  handshake chan *events.ServerMessage
  name string
  id string
}

func makeNewClient(hub *Hub, w http.ResponseWriter, r *http.Request) *Client{
  conn, err := upgrader.Upgrade(w, r, nil)
  if err != nil {
    fmt.Println("Error upgrading:", err)
    return nil
  }

  random_adjective := adjectives[rand.Intn(len(adjectives))]
  random_noun := nouns[rand.Intn(len(nouns))]

  return &Client {
    conn: conn,
    hub: hub,
		send: make(chan *events.ServerMessage),
    handshake: make(chan *events.ServerMessage),
    name: fmt.Sprintf("%s %s", random_adjective, random_noun),
    id: uuid.New().String(),
  }
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
    var msg events.ClientMessage 
    if err := proto.Unmarshal(message, &msg); err != nil {
      log.Printf("invalid message: %v", err)
      continue
    }

    // assign server ID
    msg.ServerId = &c.hub.serverId

    // marshal it back to JSON
    protoBytes, err := proto.Marshal(&msg)
    if err != nil {
      log.Printf("marshal error: %v", err)
      continue
    }

    // publish to Redis
    c.hub.redis.Publish(c.hub.ctx, "mouse_events", protoBytes)

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

      protoBytes , _ := proto.Marshal(message)
      err := c.conn.WriteMessage(websocket.BinaryMessage, protoBytes)

      if err != nil {
        log.Println("WriteMessage error:", err)
        continue
      }
    case handshake := <-c.handshake:
      if (handshake == nil) {
        continue
      }
      w,_ := c.conn.NextWriter(websocket.BinaryMessage)
      protoBytes, _ := proto.Marshal(handshake)
      w.Write(protoBytes)
      w.Close()
    }
  }
}
