package repository

import (
	"database/sql"
	"hackathon/model"
)

type UserRepository interface {
	FindById(id string) (*model.User, error)
	FindByUserName(userName string) (*model.User, error)
	Insert(tx *sql.Tx, user *model.User) error
	Update(tx *sql.Tx, user *model.User) error
	Delete(tx *sql.Tx, id string) error
}
