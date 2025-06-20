package repository

import (
	"context"
	"hackathon/model"
)

type TrendsRepository interface {
	InsertTrends(ctx context.Context, dbtx DBTX, trends []model.Trend) error
	GetTopTrendsSince(ctx context.Context, dbtx DBTX, sinceHours int, limit int) ([]model.TrendSummary, error)
}
