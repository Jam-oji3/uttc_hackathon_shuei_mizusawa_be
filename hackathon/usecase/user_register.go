package usecase

import (
	"db/dao"
	"db/model"

	"github.com/oklog/ulid/v2"
)

type RegisterUserUseCase struct {
	UserDAO *dao.UserDAO
}

func (uc *RegisterUserUseCase) Execute(name string, age int) (string, error) {
	id := ulid.Make().String()
	user := model.User{
		Id:   id,
		Name: name,
		Age:  age,
	}
	err := uc.UserDAO.Insert(user)
	return id, err
}
