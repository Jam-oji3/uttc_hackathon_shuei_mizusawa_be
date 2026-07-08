package repository

import (
	"context"
	"hackathon/model"
)

type RepostsRepository interface {
	InsertRepost(ctx context.Context, dbtx DBTX, repost *model.Repost) error
	DeleteRepost(ctx context.Context, dbtx DBTX, userID string, postID string) error
}
