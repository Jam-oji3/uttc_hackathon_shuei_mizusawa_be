package usecase

import (
	"context"
	"database/sql"
	"hackathon/model"
	"hackathon/repository"
)

type PostGetRecentUseCase struct {
	PostRepo repository.PostsRepository
	DB       *sql.DB
}

func NewPostGetRecentUseCase(postRepo repository.PostsRepository, db *sql.DB) *PostGetRecentUseCase {
	return &PostGetRecentUseCase{
		PostRepo: postRepo,
		DB:       db,
	}
}

func (uc *PostGetRecentUseCase) Execute(ctx context.Context, userId string, limit, offset int) (*[]model.PostWithUserAndCounts, error) {
	posts, err := uc.PostRepo.FindAllWithCounts(ctx, uc.DB, userId, limit, offset)
	if err != nil {
		return nil, err
	}
	return posts, nil
}
