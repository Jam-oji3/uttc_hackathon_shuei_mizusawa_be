package usecase

import (
	"context"
	"database/sql"
	"hackathon/repository"
)

type LikeDeleteUseCase struct {
	TxExecutor repository.TransactionExecutor
	LikeRepo   repository.LikesRepository
	DB         *sql.DB
}

func NewLikeDeleteUseCase(txExecutor repository.TransactionExecutor, likeRepo repository.LikesRepository, db *sql.DB) *LikeDeleteUseCase {
	return &LikeDeleteUseCase{
		TxExecutor: txExecutor,
		LikeRepo:   likeRepo,
		DB:         db,
	}
}

func (uc *LikeDeleteUseCase) Execute(ctx context.Context, userId, postId string) error {
	_, err := uc.TxExecutor.DoInTx(ctx, uc.DB, func(ctx context.Context, tx *sql.Tx) (interface{}, error) {
		if err := uc.LikeRepo.DeleteLike(ctx, tx, userId, postId); err != nil {
			return nil, err
		}
		return nil, nil
	})
	if err != nil {
		return err
	}
	return nil
}
