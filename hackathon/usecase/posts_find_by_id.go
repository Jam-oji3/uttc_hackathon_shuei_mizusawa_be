package usecase

import (
	"context"
	"database/sql"
	"hackathon/model"
	"hackathon/repository"
)

type PostFindByIdUseCase struct {
	PostRepo repository.PostsRepository
	DB       *sql.DB
}

func NewPostFindByIdUseCase(postRepo repository.PostsRepository, db *sql.DB) *PostFindByIdUseCase {
	return &PostFindByIdUseCase{
		PostRepo: postRepo,
		DB:       db,
	}
}

func (uc *PostFindByIdUseCase) Execute(ctx context.Context, userId, postId string) (*model.PostWithUserAndCounts, error) {
	post, err := uc.PostRepo.FindPostWithStatsById(ctx, uc.DB, userId, postId)
	if err != nil {
		return nil, err
	}
	return post, nil
}
