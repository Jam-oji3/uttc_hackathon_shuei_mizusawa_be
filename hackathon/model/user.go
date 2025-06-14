package model

import (
	"database/sql"
	"time"
)

type User struct {
	Id          string
	UserName    string
	DisplayName string
	Email       string
	Bio         sql.NullString
	IconURL     sql.NullString
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
