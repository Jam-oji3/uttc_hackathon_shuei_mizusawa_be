package usecase

import (
	"context"
	"database/sql"
	"errors"
	"hackathon/model"
	"hackathon/repository"
	"hackathon/util"
	"time"
)

type LikeCreateUseCase struct {
	TxExecutor repository.TransactionExecutor
	LikeRepo   repository.LikesRepository
	DB         *sql.DB
}

func NewLikeCreateUseCase(txExecutor repository.TransactionExecutor, likeRepo repository.LikesRepository, db *sql.DB) *LikeCreateUseCase {
	return &LikeCreateUseCase{
		TxExecutor: txExecutor,
		LikeRepo:   likeRepo,
		DB:         db,
	}
}

func (uc *LikeCreateUseCase) Execute(ctx context.Context, userId, postId string) (*model.Like, error) {
	id := util.GenerateULID()
	now := time.Now()

	if userId == "" || postId == "" {
		return nil, errors.New("invalid user id or post id")
	}

	like := model.Like{
		Id:        id,
		UserId:    userId,
		PostId:    postId,
		CreatedAt: now,
	}

	_, err := uc.TxExecutor.DoInTx(ctx, uc.DB, func(ctx context.Context, tx *sql.Tx) (interface{}, error) {
		if err := uc.LikeRepo.InsertLike(ctx, tx, &like); err != nil {
			return nil, err
		}
		return nil, nil
	})

	if err != nil {
		return nil, err
	}

	return &like, nil
}
