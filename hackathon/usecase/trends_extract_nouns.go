package usecase

import (
	"context"
	"database/sql"
	"hackathon/model"
	"hackathon/repository"
	"hackathon/util"
	"time"
)

type TrendExtractNounsUseCase struct {
	TxExecutor repository.TransactionExecutor
	TrendRepo  repository.TrendsRepository
	DB         *sql.DB
}

func NewTrendExtractNounsUseCase(txExecutor repository.TransactionExecutor, TrendRepo repository.TrendsRepository, DB *sql.DB) *TrendExtractNounsUseCase {
	return &TrendExtractNounsUseCase{
		TxExecutor: txExecutor,
		TrendRepo:  TrendRepo,
		DB:         DB,
	}
}

func (uc *TrendExtractNounsUseCase) Execute(ctx context.Context, text string) error {
	nouns, err := util.ExtractNouns(text)
	if err != nil {
		return err
	}
	// 重複排除
	unique := make(map[string]struct{})
	for _, noun := range nouns {
		unique[noun] = struct{}{}
	}
	//時間単位で切り捨て
	hour := time.Now().Truncate(time.Hour)

	var trends []model.Trend
	for word := range unique {
		trends = append(trends, model.Trend{
			Id:    util.GenerateULID(), // IDはユニークなULID
			Word:  word,
			Count: 1,
			Hour:  hour,
		})
	}

	_, err = uc.TxExecutor.DoInTx(ctx, uc.DB, func(ctx context.Context, tx *sql.Tx) (interface{}, error) {
		if err := uc.TrendRepo.InsertTrends(ctx, tx, trends); err != nil {
			return nil, err
		}
		return nil, nil
	})

	if err != nil {
		return err
	}

	return nil
}
