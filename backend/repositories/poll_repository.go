package repositories

import (
	"context"

	"live-polling-backend/models"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type PollRepository struct {
	collection *mongo.Collection
}

func NewPollRepository(collection *mongo.Collection) *PollRepository {
	return &PollRepository{
		collection: collection,
	}
}

func (r *PollRepository) CreatePoll(
	ctx context.Context,
	poll *models.Poll,
) error {

	_, err := r.collection.InsertOne(ctx, poll)

	return err
}

func (r *PollRepository) GetPollByID(
	ctx context.Context,
	id bson.ObjectID,
) (*models.Poll, error) {

	var poll models.Poll

	err := r.collection.FindOne(
		ctx,
		bson.M{"_id": id},
	).Decode(&poll)

	if err != nil {
		return nil, err
	}

	return &poll, nil
}

func (r *PollRepository) ClosePoll(
	ctx context.Context,
	id bson.ObjectID,
) error {

	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": id},
		bson.M{
			"$set": bson.M{
				"isActive": false,
			},
		},
	)

	return err
}
