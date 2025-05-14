package usecase

import (
	"db/dao"
	"db/model"
)

type SearchUserUseCase struct {
	UserDAO *dao.UserDAO
}

func (uc *SearchUserUseCase) Execute(name string) ([]model.User, error) {
	return uc.UserDAO.FindByName(name)
}
