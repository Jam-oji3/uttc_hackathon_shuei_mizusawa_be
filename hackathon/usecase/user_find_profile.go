package usecase

import (
	"context"
	"database/sql"
	"hackathon/model"
	"hackathon/repository"
)

type UserFindProfileUseCase struct {
	UserRepo repository.UserRepository
	DB       *sql.DB
}

func NewUserFindProfileUseCase(userRepo repository.UserRepository, db *sql.DB) *UserFindProfileUseCase {
	return &UserFindProfileUseCase{
		UserRepo: userRepo,
		DB:       db,
	}
}

func (uc *UserFindProfileUseCase) Execute(ctx context.Context, username string) (*model.UserProfile, error) {
	prof, err := uc.UserRepo.FindProfileByUsername(ctx, uc.DB, username)
	if err != nil {
		return nil, err
	}
	return prof, nil
}
