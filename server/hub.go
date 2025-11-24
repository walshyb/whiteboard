package main

import (
  "github.com/redis/go-redis/v9"
  "encoding/json"
  "log"
  "context"
)

type Hub struct{
  clients map[*Client]bool
  register chan *Client
  unregister chan *Client
  broadcast chan *InboundEvent
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

        // Remove client from active list of connections
        _, err := hub.redis.Del(hub.ctx, client.id).Result()
        if err != nil {
          log.Fatalf("Error deleting client: %v", err)
        }

        payload := OutboundEvent{
          Type: "client_disconnect",
          ClientName: client.name,
        }

        for c := range hub.clients {
          select {
          case c.send <- &payload:
          default:
            close(c.send)
            delete(hub.clients, c)
          }
        }
      }
    case message := <-hub.broadcast:
      if message.ClientId == "" {
        continue
      }

      clientName, err := hub.redis.Get(hub.ctx, message.ClientId).Result()

      if err != nil {
        println("error getting client ID from redis")
      }

      coords, _ := message.Data.(Coordinates)
      payload := OutboundEvent{
        ClientName: clientName,
        Data: coords,
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
      var m InboundEvent

      //if m.ServerId == hub.serverId {
      //  continue
      //}

      if err := json.Unmarshal([]byte(msg.Payload), &m); err != nil {
        println("Error reading JSON inbound message")
        continue
      }

      // Broadcast to our local clients
      hub.broadcast <- &m
    }
  }()
}
