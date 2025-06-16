package model

import (
	"time"
)

type User struct {
	Id          string    `json:"id"`
	UserName    string    `json:"username"`
	DisplayName string    `json:"displayName"`
	Email       string    `json:"email"`
	Bio         string    `json:"bio"`
	IconURL     string    `json:"iconUrl"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}
