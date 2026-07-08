package repository

import (
	"context"
	"hackathon/model"
)

type FollowsRepository interface {
	InsertFollow(ctx context.Context, dbtx DBTX, follow *model.Follow) error
	DeleteFollow(ctx context.Context, dbtx DBTX, followerId, followedId string) error
}
