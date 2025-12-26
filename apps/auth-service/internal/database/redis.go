package database

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

// NewRedisClient creates a new Redis client and connects to the server
func NewRedisClient(host string, port string, password string) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", host, port),
		Password: password,
		DB:       0, // use default DB
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Ping the server to verify connection
	if _, err := client.Ping(ctx).Result(); err != nil {
		return nil, err
	}

	return client, nil
}