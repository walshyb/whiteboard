package main

import (
    "fmt"
    "github.com/gorilla/websocket"
    "net/http"
    "sync"
)

var upgrader = websocket.Upgrader{
  ReadBufferSize:  1024,
	WriteBufferSize: 1024,
    CheckOrigin: func(r *http.Request) bool {
      return true;
    },
}

var clients = make(map[*websocket.Conn]bool) // Connected clients
var broadcast = make(chan []byte)            // Broadcast channel
var mutex = &sync.Mutex{}                    // Protect clients map

func wsHandler(hub *Hub, w http.ResponseWriter, r *http.Request) {
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
       fmt.Println("Error upgrading:", err)
       return
    }

    client := &Client {
      conn: conn,
      hub: hub,
      send: make(chan *Message),
    }

    client.hub.register <- client

    go client.writePump()
    go client.readPump()
}

func main() {
    hub := &Hub {
      clients: make(map[*Client]bool),
      register: make(chan *Client),
      unregister: make(chan *Client),
      broadcast: make(chan *Message),
    }
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
