package usecase

import (
	"database/sql"
	"hackathon/dao"
	"hackathon/model"
)

type UserUpdateUseCase struct {
	UserDAO *dao.UserDAO
	DB      *sql.DB
}

func (uc *UserUpdateUseCase) Execute(user *model.User) error {
	_, txErr := dao.DoInTx(uc.DB, func(tx *sql.Tx) (interface{}, error) {
		if err := uc.UserDAO.Update(tx, user); err != nil {
			return nil, err
		}
		return nil, nil
	})
	return txErr
}
