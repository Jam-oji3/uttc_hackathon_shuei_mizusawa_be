package repository

import (
	"context"
	"hackathon/model"
)

type LikesRepository interface {
	InsertLike(ctx context.Context, dbtx DBTX, like *model.Like) error
	DeleteLike(ctx context.Context, dbtx DBTX, userID string, postID string) error
}
