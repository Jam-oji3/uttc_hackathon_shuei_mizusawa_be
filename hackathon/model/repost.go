package model

import "time"

type Repost struct {
	Id        string    `json:"id"`
	UserId    string    `json:"userId"`
	PostId    string    `json:"postId"`
	CreatedAt time.Time `json:"createdAt"`
}
