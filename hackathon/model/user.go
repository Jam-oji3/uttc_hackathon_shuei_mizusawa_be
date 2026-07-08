package model

import (
	"time"
)

type User struct {
	Id          string    `json:"id"`
	Username    string    `json:"username"`
	DisplayName string    `json:"displayName"`
	Email       string    `json:"email"`
	Bio         string    `json:"bio"`
	IconURL     string    `json:"iconUrl"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type UserProfile struct {
	Id          string    `json:"id"`
	Username    string    `json:"username"`
	DisplayName string    `json:"displayName"`
	Bio         string    `json:"bio"`
	IconURL     string    `json:"iconUrl"`
	CreatedAt   time.Time `json:"createdAt"`
	Stats       struct {
		FollowingCount int64 `json:"followingCount"`
		FollowerCount  int64 `json:"followerCount"`
		PostCount      int64 `json:"postCount"`
	} `json:"stats"`
	IsFollowing bool `json:"isFollowing"`
}
