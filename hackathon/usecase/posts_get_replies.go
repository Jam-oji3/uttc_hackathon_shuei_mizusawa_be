package usecase

import (
	"context"
	"database/sql"
	"hackathon/model"
	"hackathon/repository"
)

type PostGetRepliesUseCase struct {
	PostRepo repository.PostsRepository
	DB       *sql.DB
}

func NewPostGetRepliesUseCase(postRepo repository.PostsRepository, db *sql.DB) *PostGetRepliesUseCase {
	return &PostGetRepliesUseCase{
		PostRepo: postRepo,
		DB:       db,
	}
}

func (uc *PostGetRepliesUseCase) Execute(ctx context.Context, userId, parentPostId string, limit, offset int) (*[]model.PostWithUserAndCounts, error) {
	posts, err := uc.PostRepo.FindRepliesWithStats(ctx, uc.DB, userId, parentPostId, limit, offset)
	if err != nil {
		return nil, err
	}
	return posts, nil
}
