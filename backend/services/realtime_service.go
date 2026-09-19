package services

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type RealtimeService struct {
	redisClient      *redis.Client
	webSocketService *WebSocketService
}

func NewRealtimeService(
	redisClient *redis.Client,
	webSocketService *WebSocketService,
) *RealtimeService {
	return &RealtimeService{
		redisClient:      redisClient,
		webSocketService: webSocketService,
	}
}

func (s *RealtimeService) ListenToPoll(
	ctx context.Context,
	pollID string,
) {
	channel := fmt.Sprintf("poll:%s", pollID)

	pubsub := s.redisClient.Subscribe(ctx, channel)
	defer pubsub.Close()

	fmt.Println("Listening for realtime updates on:", channel)

	for {
		message, err := pubsub.ReceiveMessage(ctx)
		if err != nil {
			fmt.Println("Realtime listener stopped:", err)
			return
		}

		s.webSocketService.Broadcast(
			pollID,
			[]byte(message.Payload),
		)
	}
}
