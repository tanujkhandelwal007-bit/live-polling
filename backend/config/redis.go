package config

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

func ConnectRedis() (*redis.Client, error) {

	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	ctx := context.Background()

	err := client.Ping(ctx).Err()
	if err != nil {
		return nil, err
	}

	fmt.Println("Redis connected successfully!")

	return client, nil
}
