package repository

import (
	"context"
	"database/sql"
	"hackathon/model"
)

type UserRepository interface {
	FindById(ctx context.Context, id string) (*model.User, error)
	FindByUserName(ctx context.Context, userName string) (*model.User, error)
	Insert(ctx context.Context, tx *sql.Tx, user *model.User) error
	Update(ctx context.Context, tx *sql.Tx, user *model.User) error
	Delete(ctx context.Context, tx *sql.Tx, id string) error
}
