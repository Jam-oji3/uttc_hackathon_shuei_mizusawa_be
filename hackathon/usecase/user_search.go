package usecase

import (
	"context"
	"hackathon/model"
	"hackathon/repository"
)

type UserSearchUseCase struct {
	UserRepo repository.UserRepository
}

func (uc *UserSearchUseCase) Execute(ctx context.Context, userName string) (*model.User, error) {
	user, err := uc.UserRepo.FindByUserName(ctx, userName)
	if err != nil {
		return nil, err
	}
	return user, nil
}
