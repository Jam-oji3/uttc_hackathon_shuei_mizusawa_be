package repository

import (
	"context"
	"hackathon/model"
)

type NotificationsRepository interface {
	GetNotifications(ctx context.Context, dbtx DBTX, userId string, limit int) ([]*model.Notification, error)
}
