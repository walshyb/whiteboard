package main

import (
  "fmt"
  "github.com/gorilla/websocket"
  "net/http"
  "math/rand"
  "github.com/redis/go-redis/v9"
  "strconv"
  "time"
  "github.com/google/uuid"
  "context"
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
  conn, err := upgrader.Upgrade(w, r, nil)
  if err != nil {
    fmt.Println("Error upgrading:", err)
    return
  }

  random_adjective := adjectives[rand.Intn(len(adjectives))]
  random_noun := nouns[rand.Intn(len(nouns))]

  client := &Client {
    conn: conn,
    hub: hub,
    send: make(chan *OutboundMessage),
    handshake: make(chan *Handshake),
    name: fmt.Sprintf("%s %s", random_adjective, random_noun),
    id: uuid.New().String(),
  }

  client.hub.register <- client

  go client.writePump()
  go client.readPump()
}

func main() {
  rdb := redis.NewClient(&redis.Options{
    Addr:     "localhost:6379",
    Password: "", // no password set
    DB:       0,  // use default DB
  })
  hub := &Hub {
    clients: make(map[*Client]bool),
    register: make(chan *Client),
    unregister: make(chan *Client),
    broadcast: make(chan *InboundMessage),
    redis: rdb,
    serverId: serverName,
    ctx: context.Background(),
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
