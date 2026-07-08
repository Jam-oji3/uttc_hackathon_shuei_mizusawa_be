package usecase

import (
	"context"
	"database/sql"
	"hackathon/model"
	"hackathon/repository"
)

type PostGetByUserUseCase struct {
	PostRepo repository.PostsRepository
	DB       *sql.DB
}

func NewPostGetByUserUseCase(postRepo repository.PostsRepository, db *sql.DB) *PostGetByUserUseCase {
	return &PostGetByUserUseCase{
		PostRepo: postRepo,
		DB:       db,
	}
}

func (uc *PostGetByUserUseCase) Execute(ctx context.Context, targetUserId, viewerUserId string, limit, offset int) (*[]model.PostWithUserAndCounts, error) {
	posts, err := uc.PostRepo.FindPostsByUserIdWithStats(ctx, uc.DB, targetUserId, viewerUserId, limit, offset)
	if err != nil {
		return nil, err
	}
	return posts, nil
}
