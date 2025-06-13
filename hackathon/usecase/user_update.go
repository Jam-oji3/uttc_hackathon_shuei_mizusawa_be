package usecase

import (
	"context"
	"database/sql"
	"hackathon/model"
	"hackathon/repository"
)

type UserUpdateUseCase struct {
	TxExecutor repository.TransactionExecutor
	UserRepo   repository.UserRepository
	DB         *sql.DB
}

func (uc *UserUpdateUseCase) Execute(ctx context.Context, user *model.User) error {
	_, txErr := uc.TxExecutor.DoInTx(ctx, uc.DB, func(ctx context.Context, tx *sql.Tx) (interface{}, error) {
		if err := uc.UserRepo.Update(ctx, tx, user); err != nil {
			return nil, err
		}
		return nil, nil
	})
	return txErr
}
