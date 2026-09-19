package services

import (
	"context"
	"errors"
	"strings"
	"time"

	"live-polling-backend/models"
	"live-polling-backend/repositories"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type PollService struct {
	repository *repositories.PollRepository
}

func NewPollService(repository *repositories.PollRepository) *PollService {
	return &PollService{
		repository: repository,
	}
}

func (s *PollService) CreatePoll(
	ctx context.Context,
	question string,
	options []string,
	createdBy string,
) (*models.Poll, error) {

	question = strings.TrimSpace(question)

	if question == "" {
		return nil, errors.New("poll question cannot be empty")
	}

	if len(options) < 2 {
		return nil, errors.New("poll must have at least 2 options")
	}

	for i := range options {
		options[i] = strings.TrimSpace(options[i])

		if options[i] == "" {
			return nil, errors.New("poll options cannot be empty")
		}
	}

	poll := &models.Poll{
		ID:        bson.NewObjectID(),
		Question:  question,
		Options:   options,
		CreatedBy: createdBy,
		IsActive:  true,
		CreatedAt: time.Now(),
	}

	err := s.repository.CreatePoll(ctx, poll)
	if err != nil {
		return nil, err
	}

	return poll, nil
}

func (s *PollService) GetPollByID(
	ctx context.Context,
	id bson.ObjectID,
) (*models.Poll, error) {

	return s.repository.GetPollByID(ctx, id)
}

func (s *PollService) ClosePoll(
	ctx context.Context,
	pollID bson.ObjectID,
	userID string,
) error {

	poll, err := s.repository.GetPollByID(
		ctx,
		pollID,
	)

	if err != nil {
		return errors.New("poll not found")
	}

	if poll.CreatedBy != userID {
		return errors.New("you are not allowed to close this poll")
	}

	if !poll.IsActive {
		return errors.New("poll is already closed")
	}

	err = s.repository.ClosePoll(
		ctx,
		pollID,
	)

	if err != nil {
		return err
	}

	return nil
}
