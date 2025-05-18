package usecase

import (
	"database/sql"
	"hackathon/dao"
	"hackathon/model"
	"time"

	"github.com/oklog/ulid/v2"
)

type UserRegisterUseCase struct {
	UserDAO *dao.UserDAO
	DB      *sql.DB
}

func (uc *UserRegisterUseCase) Execute(userName string, displayName string, bio string, iconURL string, email string) (string, error) {
	id := ulid.Make().String()
	now := time.Now()
	user := model.User{
		Id:          id,
		UserName:    userName,
		DisplayName: displayName,
		Email:       email,
		Bio:         bio,
		IconURL:     iconURL,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	_, txErr := dao.DoInTx(uc.DB, func(tx *sql.Tx) (interface{}, error) {
		if err := uc.UserDAO.Insert(tx, &user); err != nil {
			return nil, err
		}
		return nil, nil
	})
	if txErr != nil {
		return "", txErr
	}
	return id, nil
}
