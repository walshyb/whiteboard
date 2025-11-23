package main

type Hub struct{
  clients map[*Client]bool
  register chan *Client
  unregister chan *Client
  broadcast chan *Message
}

func (hub *Hub) run() {
  for {
    select {
    case client := <-hub.register: 
      hub.clients[client] = true
    case client := <- hub.unregister:
      if _, ok := hub.clients[client]; ok {
        delete(hub.clients, client)
        close(client.send)
      }
    case message := <- hub.broadcast:
      for client := range hub.clients {
        if client == message.client {
          continue
        }

        select {
        case client.send <- message:
        default:
          close(client.send)
          delete(hub.clients, client)
        }
      }
    }
  }
}

