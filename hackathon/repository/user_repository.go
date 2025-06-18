package repository

import (
	"context"
	"hackathon/model"
)

type UserRepository interface {
	FindById(ctx context.Context, dbtx DBTX, id string) (*model.User, error)
	FindByUserName(ctx context.Context, dbtx DBTX, userName string) (*model.User, error)
	FindProfileByUsername(ctx context.Context, dbtx DBTX, username string) (*model.UserProfile, error)
	Insert(ctx context.Context, dbtx DBTX, user *model.User) error
	Update(ctx context.Context, dbtx DBTX, user *model.User) error
	Delete(ctx context.Context, dbtx DBTX, id string) error
}
