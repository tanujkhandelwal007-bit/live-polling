package repositories

import (
	"context"

	"live-polling-backend/models"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

type VoteRepository struct {
	collection *mongo.Collection
}

func NewVoteRepository(collection *mongo.Collection) *VoteRepository {
	return &VoteRepository{
		collection: collection,
	}
}

func (r *VoteRepository) CreateVote(
	ctx context.Context,
	vote *models.Vote,
) error {

	_, err := r.collection.InsertOne(ctx, vote)

	return err
}
