package usecase

import (
	"hackathon/dao"
	"hackathon/model"
)

type UserSearchUseCase struct {
	UserDAO *dao.UserDAO
}

func (uc *UserSearchUseCase) Execute(username string) (*model.User, error) {
	user, err := uc.UserDAO.FindByUserName(username)
	if err != nil {
		return nil, err
	}
	return user, nil
}
