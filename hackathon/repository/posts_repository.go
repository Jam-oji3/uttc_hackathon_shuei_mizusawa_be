package repository

import (
	"context"
	"hackathon/model"
)

type PostsRepository interface {
	InsertPost(ctx context.Context, dbtx DBTX, post model.Post) error
	DeletePost(ctx context.Context, dbtx DBTX, id string) error
	FindPostWithStatsById(ctx context.Context, dbtx DBTX, userId string, postId string) (*model.PostWithUserAndCounts, error)
	FindPostsWithStats(ctx context.Context, dbtx DBTX, userId string, limit, offset int) (*[]model.PostWithUserAndCounts, error)
	FindRepliesWithStats(ctx context.Context, dbtx DBTX, userId string, parentPostId string, limit int, offset int) (*[]model.PostWithUserAndCounts, error)
	FindPostsByUserIdWithStats(ctx context.Context, dbtx DBTX, targetUserId string, viewerUserId string, limit int, offset int) (*[]model.PostWithUserAndCounts, error)
	SearchPostsByKeywordWithStats(ctx context.Context, dbtx DBTX, viewrUserId string, keyword string, limit int, offset int) (*[]model.PostWithUserAndCounts, error)
}
