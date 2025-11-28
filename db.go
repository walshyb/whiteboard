package main 

import (
	"os"
	"time"
	"context"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func makeRedisClient() *redis.Client {
	const defaultRedisAddr = "localhost:6379" 
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = defaultRedisAddr
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	return rdb
}

func makeMongoClient(ctx context.Context) *mongo.Client {
	contextWithTimeout, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	const defaultMongoURI = "mongodb://localhost:27017/whiteboard" 
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		mongoURI = defaultMongoURI
	}

	client, err := mongo.Connect(options.Client().ApplyURI(mongoURI))
	if err != nil {
		//return nil, fmt.Errorf("mongo connect error: %w", err)
		println("mongo connect error")
		return nil
	}

	// Ping to verify we can actually reach the server
	if err := client.Ping(contextWithTimeout, nil); err != nil {
		//return nil, fmt.Errorf("mongo ping error: %w", err)
		println("mongo ping error")
		return nil
	}

	return client
}
