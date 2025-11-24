package main

import (
  "github.com/redis/go-redis/v9"
  "encoding/json"
  "context"
  "fmt"
)

type Hub struct{
  clients map[*Client]bool
  register chan *Client
  unregister chan *Client
  broadcast chan *InboundMessage
  redis *redis.Client
  serverId string
  ctx context.Context
}

func (hub *Hub) run() {
  for {
    select {
    case client := <-hub.register: 
      hub.clients[client] = true
      handshake := Handshake{
        ClientId: client.id,
      }
      hub.redis.Set(hub.ctx, client.id, client.name, 0).Err()
      client.handshake <- &handshake
    case client := <- hub.unregister:
      if _, ok := hub.clients[client]; ok {
        delete(hub.clients, client)
        close(client.handshake)
        close(client.send)
      }
    case message := <-hub.broadcast:
      if message.ClientId == "" {
        continue
      }

      clientName, err := hub.redis.Get(hub.ctx, message.ClientId).Result()

      if err != nil {
        println("error getting client ID from redis")
      }

      payload := OutboundMessage{
        ClientName: clientName,
        Data: message.Data,
        Type: message.Type,
      }

      for client := range hub.clients {
        if client.id == message.ClientId {
          continue
        }

        select {
        case client.send <- &payload:
        default:
          close(client.send)
          delete(hub.clients, client)
        }
      }
    }
  }
}

func (hub *Hub) subscribeRedis() {
  sub := hub.redis.Subscribe(hub.ctx, "mouse_events")
  ch := sub.Channel()

  go func() {
    for msg := range ch {
      var m InboundMessage

      jsonData, err := json.MarshalIndent(msg.Payload, "", "  ") // Use 2 spaces for indentation
      if err != nil {
        fmt.Println("Error marshalling JSON:", err)
        return
      }

      //if m.ServerId == hub.serverId {
      //  continue
      //}

      // Print the pretty-printed JSON
      fmt.Println(string(jsonData))

      if err := json.Unmarshal([]byte(msg.Payload), &m); err != nil {
        println("Error reading JSON inbound message")
        continue
      }

      // Broadcast to our local clients
      hub.broadcast <- &m
    }
  }()
}
