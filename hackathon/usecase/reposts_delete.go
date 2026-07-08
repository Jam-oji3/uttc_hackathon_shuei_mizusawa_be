package usecase

import (
	"context"
	"database/sql"
	"hackathon/repository"
)

type RepostDeleteUseCase struct {
	TxExecutor repository.TransactionExecutor
	RepostRepo repository.RepostsRepository
	DB         *sql.DB
}

func NewRepostDeleteUseCase(txExecutor repository.TransactionExecutor, repostRepo repository.RepostsRepository, db *sql.DB) *RepostDeleteUseCase {
	return &RepostDeleteUseCase{
		TxExecutor: txExecutor,
		RepostRepo: repostRepo,
		DB:         db,
	}
}

func (uc *RepostDeleteUseCase) Execute(ctx context.Context, userId, postId string) error {
	_, err := uc.TxExecutor.DoInTx(ctx, uc.DB, func(ctx context.Context, tx *sql.Tx) (interface{}, error) {
		if err := uc.RepostRepo.DeleteRepost(ctx, tx, userId, postId); err != nil {
			return nil, err
		}
		return nil, nil
	})
	if err != nil {
		return err
	}
	return nil
}
