package usecase

import (
	"context"
	"database/sql"
	"hackathon/model"
	"hackathon/repository"
	"time"

	"github.com/oklog/ulid/v2"
)

type UserRegisterUseCase struct {
	TxExecutor repository.TransactionExecutor
	UserRepo   repository.UserRepository
	DB         *sql.DB
}

func NewUserRegisterUseCase(txExecutor repository.TransactionExecutor, userRepo repository.UserRepository, db *sql.DB) *UserRegisterUseCase {
	return &UserRegisterUseCase{TxExecutor: txExecutor, UserRepo: userRepo, DB: db}
}

func (uc *UserRegisterUseCase) Execute(ctx context.Context, userName string, displayName string, bio string, iconURL string, email string) (string, error) {
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
	_, txErr := uc.TxExecutor.DoInTx(ctx, uc.DB, func(ctx context.Context, tx *sql.Tx) (interface{}, error) {
		if err := uc.UserRepo.Insert(ctx, tx, &user); err != nil {
			return nil, err
		}
		return nil, nil
	})
	if txErr != nil {
		return "", txErr
	}
	return id, nil
}
