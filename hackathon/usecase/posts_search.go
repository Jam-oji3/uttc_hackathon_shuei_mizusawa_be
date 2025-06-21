package usecase

import (
	"context"
	"database/sql"
	"hackathon/model"
	"hackathon/repository"
)

type PostSearchUseCase struct {
	PostRepo repository.PostsRepository
	DB       *sql.DB
}

func NewPostSearchUseCase(
	postRepo repository.PostsRepository,
	db *sql.DB,
) *PostSearchUseCase {
	return &PostSearchUseCase{
		PostRepo: postRepo,
		DB:       db,
	}
}

func (uc *PostSearchUseCase) SearchPostsByKeyword(
	ctx context.Context,
	viewerUserId string,
	keyword string,
	limit, offset int,
) (*[]model.PostWithUserAndCounts, error) {
	posts, err := uc.PostRepo.SearchPostsByKeywordWithStats(ctx, uc.DB, viewerUserId, keyword, limit, offset)
	if err != nil {
		return nil, err
	}
	return posts, nil
}
