package usecase

import (
	"hackathon/model"
	"hackathon/repository"
)

type UserSearchUseCase struct {
	UserRepo repository.UserRepository
}

func (uc *UserSearchUseCase) Execute(username string) (*model.User, error) {
	user, err := uc.UserRepo.FindByUserName(username)
	if err != nil {
		return nil, err
	}
	return user, nil
}
