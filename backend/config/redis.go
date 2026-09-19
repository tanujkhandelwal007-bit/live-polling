package config

import (
	"context"
	"fmt"
	"os"

	"github.com/redis/go-redis/v9"
)

func ConnectRedis() (*redis.Client, error) {

	redisURL := os.Getenv("REDIS_URL")

	// Local development fallback
	if redisURL == "" {
		redisURL = "redis://localhost:6379"
	}

	options, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, err
	}

	client := redis.NewClient(options)

	ctx := context.Background()

	err = client.Ping(ctx).Err()
	if err != nil {
		return nil, err
	}

	fmt.Println("Redis connected successfully!")

	return client, nil
}
