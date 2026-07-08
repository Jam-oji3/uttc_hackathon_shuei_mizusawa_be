package usecase

import (
	"context"
	"database/sql"
	"hackathon/model"
	"hackathon/repository"
)

type TrendGetTopUseCase struct {
	TrendRepo repository.TrendsRepository
	DB        *sql.DB
}

func NewTrendGetTopUseCase(trendRepo repository.TrendsRepository, db *sql.DB) *TrendGetTopUseCase {
	return &TrendGetTopUseCase{
		TrendRepo: trendRepo,
		DB:        db,
	}
}

func (uc *TrendGetTopUseCase) Execute(ctx context.Context, sinceHours, limit int) ([]model.TrendSummary, error) {
	trends, err := uc.TrendRepo.GetTopTrendsSince(ctx, uc.DB, sinceHours, limit)
	if err != nil {
		return nil, err
	}
	return trends, nil
}
