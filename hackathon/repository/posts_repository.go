package repository

import (
	"context"
	"hackathon/model"
)

type PostsRepository interface {
	FindPostById(ctx context.Context, dbtx DBTX, id string) (*model.Post, error)
	FindPostsByUserId(ctx context.Context, dbtx DBTX, userId string) (*[]model.Post, error)
	InsertPost(ctx context.Context, dbtx DBTX, post model.Post) error
	DeletePost(ctx context.Context, dbtx DBTX, id string) error
	FindPostsWithStats(ctx context.Context, dbtx DBTX, userId string, limit, offset int) (*[]model.PostWithUserAndCounts, error)
}
