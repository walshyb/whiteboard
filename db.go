package main 

import (
  "time"
  "context"
  "github.com/redis/go-redis/v9"
  "go.mongodb.org/mongo-driver/v2/mongo"
  "go.mongodb.org/mongo-driver/v2/mongo/options"

)

func makeRedisClient() *redis.Client {
  rdb := redis.NewClient(&redis.Options{
    Addr:     "localhost:6379",
    Password: "", // no password set
    DB:       0,  // use default DB
  })
  return rdb
}

func makeMongoClient(ctx context.Context) *mongo.Client {
  context, cancel := context.WithTimeout(ctx, 10*time.Second)
    defer cancel()

  client, err := mongo.Connect(options.Client().ApplyURI("mongodb://localhost:27017"))
    if err != nil {
        //return nil, fmt.Errorf("mongo connect error: %w", err)
        println("mongo connect error")
    }

    // Ping to verify we can actually reach the server
    if err := client.Ping(context, nil); err != nil {
        //return nil, fmt.Errorf("mongo ping error: %w", err)
        println("mongo ping error")
    }

    return client
}
