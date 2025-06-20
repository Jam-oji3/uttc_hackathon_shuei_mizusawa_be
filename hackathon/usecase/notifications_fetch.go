package usecase

import (
	"context"
	"database/sql"
	"hackathon/model"
	"hackathon/repository"
)

type NotificationFetchUseCase struct {
	NotifRepo repository.NotificationsRepository
	DB        *sql.DB
}

func NewNotificationFetchUseCase(notificationRepo repository.NotificationsRepository, db *sql.DB) *NotificationFetchUseCase {
	return &NotificationFetchUseCase{
		NotifRepo: notificationRepo,
		DB:        db,
	}
}

func (uc *NotificationFetchUseCase) Execute(ctx context.Context, userId string, limit int) ([]*model.Notification, error) {
	notifications, err := uc.NotifRepo.GetNotifications(ctx, uc.DB, userId, limit)
	if err != nil {
		return nil, err
	}
	return notifications, nil
}
