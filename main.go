package main

import (
	"os"
	"fmt"
	"net/http"
	"strconv"
	"time"
	"context"
	"strings"

	"github.com/gorilla/websocket"
	events "github.com/walshyb/whiteboard/proto"
)

var serverName = strconv.FormatInt(time.Now().UnixNano(), 10)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		origin := r.Header.Get("Origin")
		if origin == "" {
			return false
		}
		
		// Read from environment variable, default to localhost for dev
		allowedOriginsEnv := os.Getenv("ALLOWED_ORIGINS")
		currentEnv := os.Getenv("ENV")
		if currentEnv != "production" && allowedOriginsEnv == "" {
			allowedOriginsEnv = "http://localhost:5173,http://localhost:3000"
		}

		// TODO: evaluate allowedOrigins once
		
		allowedOrigins := strings.Split(allowedOriginsEnv, ",")
		for _, allowed := range allowedOrigins {
			if origin == strings.TrimSpace(allowed) {
				return true
			}
		}
		return false
	},
}

func wsHandler(hub *Hub, w http.ResponseWriter, r *http.Request) {
	// Decline connection if it's not a websocket
	if !websocket.IsWebSocketUpgrade(r) {
		http.Error(w, "WebSocket upgrade required", http.StatusUpgradeRequired)
		return
	}
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
		broadcast: make(chan *events.ClientMessage),
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

	http.HandleFunc("/health", func (w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	fmt.Println("WebSocket server started on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
