package main

import (
  "log"
  "context"
  "github.com/redis/go-redis/v9"
  "go.mongodb.org/mongo-driver/v2/mongo"
	"google.golang.org/protobuf/proto"
	 events "github.com/walshyb/whiteboard/proto"
)

type Hub struct{
  clients map[*Client]bool
  register chan *Client
  unregister chan *Client
  broadcast chan *events.ClientMessage
  redis *redis.Client
  mongo *mongo.Client
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

        payload := events.ServerMessage{
					NotificationType: &events.ServerMessage_ClientDisconnect {
						ClientDisconnect: &events.ClientDisconnectEvent{
							ClientName: client.name,
						},
					},
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
    case clientMessage := <-hub.broadcast:
      if clientMessage.ClientId == "" {
        continue
      }

      clientName, err := hub.redis.Get(hub.ctx, clientMessage.ClientId).Result()

      if err != nil {
        println("error getting client ID from redis")
				continue
      }

      serverMessage := &events.ServerMessage{
				SenderName: clientName,
				NotificationType: &events.ServerMessage_EventData {
					EventData: clientMessage.GetEvent(),
				},
      }

      for client := range hub.clients {
        if client.id == clientMessage.GetClientId() {
          continue
        }

        select {
        case client.send <- serverMessage:
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
      var m events.ClientMessage

      //if m.ServerId == hub.serverId {
      //  continue
      //}

      if err := proto.Unmarshal([]byte(msg.Payload), &m); err != nil {
        println("Error reading JSON inbound message")
        continue
      }

      // Broadcast to our local clients
      hub.broadcast <- &m
    }
  }()
}
