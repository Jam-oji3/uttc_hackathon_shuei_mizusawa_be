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

type RepostCreateUseCase struct {
	TxExecutor repository.TransactionExecutor
	RepostRepo repository.RepostsRepository
	DB         *sql.DB
}

func NewRepostCreateUseCase(txExecutor repository.TransactionExecutor, repostRepo repository.RepostsRepository, db *sql.DB) *RepostCreateUseCase {
	return &RepostCreateUseCase{
		TxExecutor: txExecutor,
		RepostRepo: repostRepo,
		DB:         db,
	}
}

func (uc *RepostCreateUseCase) Execute(ctx context.Context, userId, postId string) (*model.Repost, error) {
	id := util.GenerateULID()
	now := time.Now()

	if userId == "" || postId == "" {
		return nil, errors.New("invalid user id or post id")
	}

	repost := model.Repost{
		Id:        id,
		UserId:    userId,
		PostId:    postId,
		CreatedAt: now,
	}

	_, err := uc.TxExecutor.DoInTx(ctx, uc.DB, func(ctx context.Context, tx *sql.Tx) (interface{}, error) {
		if err := uc.RepostRepo.InsertRepost(ctx, tx, &repost); err != nil {
			return nil, err
		}
		return nil, nil
	})

	if err != nil {
		return nil, err
	}

	return &repost, nil
}
