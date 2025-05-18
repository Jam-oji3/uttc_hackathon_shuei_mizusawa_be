package model

import "time"

type User struct {
	Id          string
	UserName    string
	DisplayName string
	Email       string
	Bio         string
	IconURL     string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
