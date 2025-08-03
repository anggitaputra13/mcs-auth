package redis

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

// Client wraps the redis.Client for easier usage
type Client struct {
	*redis.Client
}

// NewRedisClient creates a new Redis client
func NewRedisClient(url string) *Client {
	opt, err := redis.ParseURL(url)
	if err != nil {
		log.Fatalf("Failed to parse Redis URL: %v", err)
	}

	client := redis.NewClient(opt)

	// Test connection
	if _, err := client.Ping(context.Background()).Result(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	log.Println("Connected to Redis successfully")
	return &Client{client}
}
