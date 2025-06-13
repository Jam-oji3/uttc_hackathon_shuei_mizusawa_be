package usecase

import (
	"context"
	"database/sql"
	"hackathon/infra/db"
	"hackathon/model"
	"hackathon/repository"
)

type UserUpdateUseCase struct {
	UserRepo repository.UserRepository
	DB       *sql.DB
}

func (uc *UserUpdateUseCase) Execute(ctx context.Context, user *model.User) error {
	_, txErr := db.DoInTx(uc.DB, func(tx *sql.Tx) (interface{}, error) {
		if err := uc.UserRepo.Update(ctx, tx, user); err != nil {
			return nil, err
		}
		return nil, nil
	})
	return txErr
}
