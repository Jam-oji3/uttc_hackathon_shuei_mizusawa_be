package mysql

import (
	"context"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"hackathon/model"
)

func TestInsertMultiple(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %v", err)
	}
	defer db.Close()

	repo := NewTrendsRepository()

	ctx := context.Background()

	trends := []model.Trend{
		{
			Id:    "01F7Y2G0X5C3Z1A4B2C3D4E5F6",
			Word:  "testword1",
			Hour:  time.Date(2025, 6, 19, 14, 0, 0, 0, time.UTC),
			Count: 3,
		},
		{
			Id:    "01F7Y2G0X5C3Z1A4B2C3D4E5F7",
			Word:  "testword2",
			Hour:  time.Date(2025, 6, 19, 14, 0, 0, 0, time.UTC),
			Count: 5,
		},
	}

	// プレースホルダーの数とSQLを検証（正規表現で）
	expectedQuery := regexp.QuoteMeta(
		"INSERT INTO trends (id, word, hour, count) VALUES (?, ?, ?, ?),(?, ?, ?, ?) ON DUPLICATE KEY UPDATE count = count + VALUES(count)",
	)

	// モックで期待されるExecの呼び出しをセット
	mock.ExpectExec(expectedQuery).
		WithArgs(
			trends[0].Id, trends[0].Word, trends[0].Hour, trends[0].Count,
			trends[1].Id, trends[1].Word, trends[1].Hour, trends[1].Count,
		).
		WillReturnResult(sqlmock.NewResult(2, 2)) // lastInsertId, rowsAffected

	err = repo.InsertTrends(ctx, db, trends)
	if err != nil {
		t.Errorf("InsertTrends returned error: %v", err)
	}

	// 全ての期待が満たされたかチェック
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %v", err)
	}
}
