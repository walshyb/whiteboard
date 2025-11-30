package main

import (
	"os"
	"fmt"
	"net/http"
	"strconv"
	"time"
	"context"
	"strings"
	"encoding/json"

	"github.com/gorilla/websocket"
	events "github.com/walshyb/whiteboard/proto"
)

var serverName = strconv.FormatInt(time.Now().UnixNano(), 10)

// Read from environment variable, default to localhost for dev
var allowedOriginsEnv = os.Getenv("ALLOWED_ORIGINS")
var currentEnv = os.Getenv("ENV")

var allowedOrigins []string

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		origin := r.Header.Get("Origin")
		if origin == "" {
			return false
		}

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

	if currentEnv != "production" && allowedOriginsEnv == "" {
		allowedOriginsEnv = "http://localhost:5173,http://localhost:3000"
	}

	allowedOrigins = strings.Split(allowedOriginsEnv, ",")

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

	http.HandleFunc("/board", func (w http.ResponseWriter, r *http.Request) {
		// TODO: allow OPTIONS, GET method
		origin := r.Header.Get("Origin")
		isAllowed := false
		// TODO: create helper
    for _, allowedOrigin := range allowedOrigins {
        if origin == allowedOrigin {
            isAllowed = true
            break
        }
    }

		if (isAllowed) {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		}

		w.Header().Set("Content-Type", "application/json")
		// Call your GetBoard function
		board, err := hub.GetBoard()
		if err != nil {
			// If something went wrong, return 500
			// TODO: one day switch to proper logs
			println("GetBoard error: %v", err)
			http.Error(w, "failed to retrieve board", http.StatusInternalServerError)
			return
		}

		// Marshal the board to JSON
		jsonData, err := json.Marshal(board)
		if err != nil {
			println("JSON marshal error: %v", err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		// Write the JSON response
		w.WriteHeader(http.StatusOK)
		w.Write(jsonData)
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
