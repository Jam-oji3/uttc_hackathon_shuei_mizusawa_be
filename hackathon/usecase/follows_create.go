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

type FollowCreateUseCase struct {
	TxExecutor repository.TransactionExecutor
	FollowRepo repository.FollowsRepository
	DB         *sql.DB
}

func NewFollowCreateUseCase(txExecutor repository.TransactionExecutor, followRepo repository.FollowsRepository, db *sql.DB) *FollowCreateUseCase {
	return &FollowCreateUseCase{
		TxExecutor: txExecutor,
		FollowRepo: followRepo,
		DB:         db,
	}
}

func (uc *FollowCreateUseCase) Execute(ctx context.Context, followerId, followedId string) error {
	id := util.GenerateULID()
	now := time.Now()

	if followerId == "" || followedId == "" {
		return errors.New("invalid follower id or followed id")
	}

	follow := model.Follow{
		Id:         id,
		FollowerId: followerId,
		FollowedId: followedId,
		CreatedAt:  now,
	}

	_, err := uc.TxExecutor.DoInTx(ctx, uc.DB, func(ctx context.Context, tx *sql.Tx) (interface{}, error) {
		if err := uc.FollowRepo.InsertFollow(ctx, tx, &follow); err != nil {
			return nil, err
		}
		return nil, nil
	})

	if err != nil {
		return err
	}

	return nil
}
