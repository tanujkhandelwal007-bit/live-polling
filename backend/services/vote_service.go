package services

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"live-polling-backend/models"
	"live-polling-backend/repositories"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type VoteService struct {
	repository     *repositories.VoteRepository
	pollRepository *repositories.PollRepository
	redisService   *RedisService
}

func NewVoteService(
	repository *repositories.VoteRepository,
	pollRepository *repositories.PollRepository,
	redisService *RedisService,
) *VoteService {
	return &VoteService{
		repository:     repository,
		pollRepository: pollRepository,
		redisService:   redisService,
	}
}

func (s *VoteService) CreateVote(
	ctx context.Context,
	pollID bson.ObjectID,
	option string,
) (*models.Vote, error) {

	option = strings.TrimSpace(option)

	if option == "" {
		return nil, errors.New("vote option cannot be empty")
	}

	// Get poll from MongoDB
	poll, err := s.pollRepository.GetPollByID(ctx, pollID)
	if err != nil {
		return nil, errors.New("poll not found")
	}

	// Check if poll is active
	if !poll.IsActive {
		return nil, errors.New("poll is not active")
	}

	// Check if selected option exists in the poll
	optionExists := false

	for _, pollOption := range poll.Options {
		if pollOption == option {
			optionExists = true
			break
		}
	}

	if !optionExists {
		return nil, errors.New("invalid poll option")
	}

	vote := &models.Vote{
		ID:        bson.NewObjectID(),
		PollID:    pollID,
		Option:    option,
		CreatedAt: time.Now(),
	}

	// Save vote permanently in MongoDB
	err = s.repository.CreateVote(ctx, vote)
	if err != nil {
		return nil, err
	}

	// Increment live vote count in Redis
	err = s.redisService.IncrementVoteCount(
		ctx,
		pollID.Hex(),
		option,
	)

	if err != nil {
		return nil, err
	}

	// Get updated vote count from Redis
	results, err := s.redisService.GetVoteResults(
		ctx,
		pollID.Hex(),
	)

	if err != nil {
		return nil, err
	}

	count := int64(0)

	if value, ok := results[option]; ok {
		fmt.Sscan(value, &count)
	}

	// Publish realtime update through Redis Pub/Sub
	err = s.redisService.PublishVoteUpdate(
		ctx,
		pollID.Hex(),
		option,
		count,
	)

	if err != nil {
		return nil, err
	}

	return vote, nil
}
