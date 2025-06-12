package usecase

import (
	"database/sql"
	"hackathon/infra/db"
	"hackathon/model"
	"hackathon/repository"
	"time"

	"github.com/oklog/ulid/v2"
)

type UserRegisterUseCase struct {
	UserRepo repository.UserRepository
	DB       *sql.DB
}

func NewUserRegisterUseCase(userRepo repository.UserRepository, db *sql.DB) *UserRegisterUseCase {
	return &UserRegisterUseCase{UserRepo: userRepo, DB: db}
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
	_, txErr := db.DoInTx(uc.DB, func(tx *sql.Tx) (interface{}, error) {
		if err := uc.UserRepo.Insert(tx, &user); err != nil {
			return nil, err
		}
		return nil, nil
	})
	if txErr != nil {
		return "", txErr
	}
	return id, nil
}
