package services

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type RedisService struct {
	client *redis.Client
}

func NewRedisService(client *redis.Client) *RedisService {
	return &RedisService{
		client: client,
	}
}

// IncrementVoteCount increases the vote count for an option.
func (s *RedisService) IncrementVoteCount(
	ctx context.Context,
	pollID string,
	option string,
) error {

	key := fmt.Sprintf("poll:%s:results", pollID)

	return s.client.HIncrBy(ctx, key, option, 1).Err()
}

// GetVoteResults returns all current vote counts.
func (s *RedisService) GetVoteResults(
	ctx context.Context,
	pollID string,
) (map[string]string, error) {

	key := fmt.Sprintf("poll:%s:results", pollID)

	return s.client.HGetAll(ctx, key).Result()
}

// PublishVoteUpdate sends a realtime event through Redis Pub/Sub.
func (s *RedisService) PublishVoteUpdate(
	ctx context.Context,
	pollID string,
	option string,
	count int64,
) error {

	channel := fmt.Sprintf("poll:%s", pollID)

	event := map[string]interface{}{
		"pollId": pollID,
		"option": option,
		"count":  count,
	}

	data, err := json.Marshal(event)
	if err != nil {
		return err
	}

	return s.client.Publish(ctx, channel, data).Err()
}
