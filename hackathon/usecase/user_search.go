package usecase

import (
	"context"
	"database/sql"
	"hackathon/model"
	"hackathon/repository"
)

type UserSearchUseCase struct {
	UserRepo repository.UserRepository
	DB       *sql.DB
}

func (uc *UserSearchUseCase) Execute(ctx context.Context, userName string) (*model.User, error) {
	user, err := uc.UserRepo.FindByUserName(ctx, uc.DB, userName)
	if err != nil {
		return nil, err
	}
	return user, nil
}
