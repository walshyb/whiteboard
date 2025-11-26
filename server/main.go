package main

import (
  "fmt"
  "net/http"
  "strconv"
  "time"
  "context"

  "github.com/gorilla/websocket"
)

var adjectives = [8]string{"bright", "silent", "rough", "narrow", "gentle", "sharp", "steady", "fragile",}
var nouns = [8]string{"river","lantern","stone", "meadow","circuit","anchor","window","compass",}
var serverName = strconv.FormatInt(time.Now().UnixNano(), 10)

var upgrader = websocket.Upgrader{
  ReadBufferSize:  1024,
  WriteBufferSize: 1024,
  CheckOrigin: func(r *http.Request) bool {
    return true;
  },
}

func wsHandler(hub *Hub, w http.ResponseWriter, r *http.Request) {
  client := makeNewClient(hub, w, r)

  if client == nil {
    return
  }

  client.hub.register <- client

  go client.writePump()
  go client.readPump()
}

func main() {
  rdb := makeRedisClient()
  ctx := context.Background()
  mongo := makeMongoClient(ctx)

  hub := &Hub {
    clients: make(map[*Client]bool),
    register: make(chan *Client),
    unregister: make(chan *Client),
    broadcast: make(chan *InboundEvent),
    redis: rdb,
    mongo: mongo,
    serverId: serverName,
    ctx: ctx,
  }

  hub.subscribeRedis()
  go hub.run()

  http.HandleFunc("/ws", func (w http.ResponseWriter, r *http.Request) {
    wsHandler(hub, w, r)
  })

  fmt.Println("WebSocket server started on :8080")
  err := http.ListenAndServe(":8080", nil)
  if err != nil {
    fmt.Println("Error starting server:", err)
  }
}
