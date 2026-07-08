package mysql

import (
	"context"
	"hackathon/model"
	"hackathon/repository"
	"strings"
)

type TrendsRepository struct{}

var _ repository.TrendsRepository = (*TrendsRepository)(nil)

func NewTrendsRepository() *TrendsRepository {
	return &TrendsRepository{}
}

func (r *TrendsRepository) InsertTrends(ctx context.Context, dbtx repository.DBTX, trends []model.Trend) error {
	if len(trends) == 0 {
		return nil
	}

	placeholders := make([]string, 0, len(trends))
	values := make([]interface{}, 0, len(trends)*4)

	for _, t := range trends {
		placeholders = append(placeholders, "(?, ?, ?, ?)")
		values = append(values, t.Id, t.Word, t.Hour, t.Count)
	}

	query := `
        INSERT INTO trends (id, word, hour, count)
        VALUES ` + strings.Join(placeholders, ",") + `
        ON DUPLICATE KEY UPDATE count = count + VALUES(count)
    `

	_, err := dbtx.ExecContext(ctx, query, values...)
	return err
}

func (r *TrendsRepository) GetTopTrendsSince(ctx context.Context, dbtx repository.DBTX, sinceHours int, limit int) ([]model.TrendSummary, error) {
	query := `
        SELECT 
            word,
            SUM(count) as total_count
        FROM trends
        WHERE hour >= NOW() - INTERVAL ? HOUR
        GROUP BY word
        ORDER BY total_count DESC
        LIMIT ?
    `

	rows, err := dbtx.QueryContext(ctx, query, sinceHours, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var trends []model.TrendSummary
	for rows.Next() {
		var t model.TrendSummary
		if err := rows.Scan(&t.Word, &t.Count); err != nil {
			return nil, err
		}
		trends = append(trends, t)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return trends, nil
}
