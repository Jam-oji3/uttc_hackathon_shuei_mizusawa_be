package usecase

import (
	"context"
	"database/sql"
	"hackathon/repository"
)

type FollowDeleteUseCase struct {
	TxExecutor repository.TransactionExecutor
	FollowRepo repository.FollowsRepository
	DB         *sql.DB
}

func NewFollowDeleteUseCase(txExecutor repository.TransactionExecutor, followRepo repository.FollowsRepository, db *sql.DB) *FollowDeleteUseCase {
	return &FollowDeleteUseCase{
		TxExecutor: txExecutor,
		FollowRepo: followRepo,
		DB:         db,
	}
}

func (uc *FollowDeleteUseCase) Execute(ctx context.Context, followerId, followedId string) error {
	_, err := uc.TxExecutor.DoInTx(ctx, uc.DB, func(ctx context.Context, tx *sql.Tx) (interface{}, error) {
		if err := uc.FollowRepo.DeleteFollow(ctx, tx, followerId, followedId); err != nil {
			return nil, err
		}
		return nil, nil
	})
	if err != nil {
		return err
	}
	return nil
}
