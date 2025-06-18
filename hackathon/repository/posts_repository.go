package repository

import (
	"context"
	"hackathon/model"
)

type PostsRepository interface {
	FindPostWithStatsById(ctx context.Context, dbtx DBTX, userId string, postId string) (*model.PostWithUserAndCounts, error)
	FindPostsByUserId(ctx context.Context, dbtx DBTX, userId string) (*[]model.Post, error)
	InsertPost(ctx context.Context, dbtx DBTX, post model.Post) error
	DeletePost(ctx context.Context, dbtx DBTX, id string) error
	FindPostsWithStats(ctx context.Context, dbtx DBTX, userId string, limit, offset int) (*[]model.PostWithUserAndCounts, error)
	FindRepliesWithStats(ctx context.Context, dbtx DBTX, userId string, parentPostId string, limit int, offset int) (*[]model.PostWithUserAndCounts, error)
}
